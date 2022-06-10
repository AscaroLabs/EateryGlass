package storage

// Данный пакет отвечает за хранение данных
// и взаимодействия с базой данных

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/AscaroLabs/EateryGlass/internal/config"
	"github.com/AscaroLabs/EateryGlass/pkg/structures"
	"github.com/lib/pq"
)

// NewDB создает новый экземпляр sql.DB
// А так же создает таблицы и заполняет их тестовыми данными
func NewDB(cfg *config.Config) *sql.DB {

	// Стока для соединения с БД
	// e.g. $ psql postgresql://dbmaster:5433/mydb?sslmode=require
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DB_user,
		cfg.DB_password,
		cfg.DB_host,
		cfg.DB_port,
		cfg.DB_name)

	// Открываем БД с через драйвер "postgres"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Проверяем, можем ли мы обращаться к открытой БД
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	// Заполняем БД тестовыми данными
	FillDB(cfg, db)

	return db
}

// FillDB заполняет БД содержимым
func FillDB(cfg *config.Config, db *sql.DB) {

	// Создаем таблицы
	ExecFromFile(cfg, db, "create_model.sql")

	// Заполняем таблицы с ресторанами и столиками
	ExecFromFile(cfg, db, "fill_db.sql")

	// Добавляем тестовых клиентов
	ExecFromFile(cfg, db, "insert_clients.sql")

	// Добавляем тестовые брони
	ExecFromFile(cfg, db, "insert_reservations.sql")
}

// ExecFromFile исполняет команды из файла file_name,
// хранящугося в директории /pkg/database/postgresql/storage/
func ExecFromFile(cfg *config.Config, db *sql.DB, file_name string) error {

	// Контекст, используемый по умолчанию в go
	ctx := context.Background()

	// Открываем файл с командами
	file, err := os.Open(
		fmt.Sprintf("%s/pkg/database/postgresql/storage/%s",
			cfg.Main_dir,
			file_name))
	if err != nil {
		return err
	}

	// Считываем данные из файла в буффер (предполагается, что файл с командами весит меньше 1MB)
	data := make([]byte, 1024)
	n, err := file.Read(data)

	// Парсим из данных команды
	queries := strings.Split((string(data[:n])), "\n")

	// Начинаем транзакцию
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// По порядку исполняем команды
	for _, q := range queries {

		// на всякий случай проверяем на пустоту,
		// врдуг случйно в файле оказались лишние переносы
		if q == "" {
			continue
		}

		// Выполняем команду в рамках транзакции
		_, execErr := tx.Exec(q)
		if execErr != nil {
			// Если врдуг что-то пошло не так, то нужно откатить всю транзакцию
			_ = tx.Rollback()
			return execErr
		}
	}

	// Если все прошло хорошо, то подтверждаем транзакцию
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// GetRestaurants получает из БД данные о ресторанах
func GetRestaurants(db *sql.DB) ([]structures.Restaurant, error) {
	q := `
		SELECT * FROM restaurants;
	`
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var restaurants []structures.Restaurant
	for rows.Next() {
		var restaurant structures.Restaurant
		if err := rows.Scan(&restaurant.ID, &restaurant.Name,
			&restaurant.Avg_time, &restaurant.Avg_price); err != nil {
			return nil, err
		}
		restaurants = append(restaurants, restaurant)
	}
	return restaurants, nil
}

// GetTablesByTime получает из БД столики, которые свободны на момент t
func GetTablesByTime(db *sql.DB, t time.Time) ([]structures.Table, error) {
	timeStamp := pq.FormatTimestamp(t)
	q := `
		SELECT DISTINCT tables.id, tables.restaurant_id, tables.capacity 
		FROM tables LEFT OUTER JOIN reservations 
		ON (tables.id = reservations.table_id) WHERE
		($1::timestamp + '2 hours'::interval <= start_time) OR 
		($1::timestamp >= end_time) OR 
		(reservations.id IS NULL) 
		ORDER BY tables.id;
	`
	rows, err := db.Query(q, timeStamp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tables []structures.Table
	for rows.Next() {
		var table structures.Table
		if err := rows.Scan(&table.ID, &table.Restaurant_id, &table.Capacity); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tables, nil
}

// PostReservations добавляет новую бронь в БД
func PostReservations(db *sql.DB, rawReservation structures.RawReservation) (structures.Reservation, error) {
	ctx := context.Background()
	squery := `
		SELECT * FROM clients  WHERE
		(name=$1) AND 
		(phone_number=$2);
	`
	// Пытаемся найти клиента в БД
	var client structures.Client
	row := db.QueryRow(squery, rawReservation.Reserved_by.Name, rawReservation.Reserved_by.Phone)
	if err := row.Scan(&client.ID, &client.Name, &client.Phone); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Добавляем клиента в БД
			client, err = AddClient(db, rawReservation.Reserved_by)
			if err != nil {
				return structures.Reservation{}, err
			}
		} else {
			return structures.Reservation{}, err
		}
	}

	// Начинаем транзакцию
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return structures.Reservation{}, err
	}
	q := `
	INSERT INTO reservations (table_id,start_time,end_time,reserved_by)
	VALUES  ($1, $2::timestamp, $2::timestamp + '2 H'::interval, $3);
	`
	startTime, err := time.Parse(time.RFC3339, rawReservation.Start_time)
	_, execErr := tx.Exec(q, rawReservation.Table_id, pq.FormatTimestamp(startTime), client.ID)
	if execErr != nil {
		_ = tx.Rollback()
		return structures.Reservation{}, execErr
	}
	if err := tx.Commit(); err != nil {
		return structures.Reservation{}, err
	}

	// Теперь хотим вернуть добавленную бронь
	var addedReservation structures.Reservation
	row = db.QueryRow(`
		SELECT * FROM reservations
		WHERE (table_id=$1) 
		AND (start_time=$2::timestamp) 
		AND (reserved_by=$3);
	`, rawReservation.Table_id, pq.FormatTimestamp(startTime), client.ID)

	if err := row.Scan(&addedReservation.ID,
		&addedReservation.Table_id,
		&addedReservation.Start_time,
		&addedReservation.End_time,
		&addedReservation.Reserved_by); err != nil {
		return structures.Reservation{}, err
	}

	return addedReservation, nil
}

// AddClient добавляет нового клиента в БД
func AddClient(db *sql.DB, rawClient structures.RawClient) (structures.Client, error) {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return structures.Client{}, err
	}
	q := `
	INSERT INTO clients (name, phone_number)
	VALUES  ($1, $2);
	`
	_, execErr := tx.Exec(q, rawClient.Name, rawClient.Phone)
	if execErr != nil {
		_ = tx.Rollback()
		return structures.Client{}, execErr
	}
	if err := tx.Commit(); err != nil {
		return structures.Client{}, err
	}

	//  Для того чтобы вернуть structures.Client не хватает только идентификатора
	//	Возьмем его из вспомогательной таблицы clients_id_seq
	q = `
		SELECT last_value FROM clients_id_seq;
	`
	row := db.QueryRow(q)
	var clientId int
	if err := row.Scan(&clientId); err != nil {
		return structures.Client{}, err
	}
	return structures.Client{
		ID:    fmt.Sprintf("%v", clientId),
		Name:  rawClient.Name,
		Phone: rawClient.Phone}, nil
}

// GetClientByRaw получет из БД данные о клиенте
// по его имени и номеру телефона
func GetClientByRaw(db *sql.DB, rawClient structures.RawClient) (structures.Client, error) {
	q := `
	SELECT * FROM clients  WHERE
	(name=$1) AND 
	(phone_number=$2);
	`
	var client structures.Client
	row := db.QueryRow(q, rawClient.Name, rawClient.Phone)
	if err := row.Scan(&client.ID, &client.Name, &client.Phone); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			client, err := AddClient(db, rawClient)
			if err != nil {
				return structures.Client{}, err
			}
			return client, nil
		}
		return structures.Client{}, err
	}
	return client, nil
}
