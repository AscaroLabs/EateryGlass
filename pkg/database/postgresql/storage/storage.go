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

// Данная функция создает новый экземпляр sql.DB
// А так же создает таблицы и заполняет их тестовыми данными
// (функция FillDB)
func NewDB(cfg *config.Config) *sql.DB {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DB_user,
		cfg.DB_password,
		cfg.DB_host,
		cfg.DB_port,
		cfg.DB_name)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	FillDB(cfg, db)

	return db
}

// Заполняем БД содержимым
func FillDB(cfg *config.Config, db *sql.DB) {
	ExecFromFile(cfg, db, "create_model.sql")
	ExecFromFile(cfg, db, "fill_db.sql")
	ExecFromFile(cfg, db, "insert_clients.sql")
	ExecFromFile(cfg, db, "insert_reservations.sql")
}

// Функция исполняет команды из файла file_name,
// хранящугося в директории /pkg/database/postgresql/storage/
func ExecFromFile(cfg *config.Config, db *sql.DB, file_name string) error {
	defer log.Println("==============================")

	ctx := context.Background()

	file, err := os.Open(
		fmt.Sprintf("%s/pkg/database/postgresql/storage/%s",
			cfg.Main_dir,
			file_name))

	log.Printf("%s is opened", file_name)

	if err != nil {
		return err
	}
	data := make([]byte, 1024)
	n, err := file.Read(data)
	queries := strings.Split((string(data[:n])), "\n")

	// for i, q := range queries {
	// 	log.Printf("Command %v --> %s \n", i, q)
	// }

	tx, err := db.BeginTx(ctx, nil)

	log.Printf("%s txn begined!\n", file_name)

	if err != nil {
		return err
	}
	for _, q := range queries {
		if q == "" {
			continue
		}

		log.Printf("start execute %s\n", q)

		_, execErr := tx.Exec(q)

		log.Printf("%s\nexecuted!", q)

		if execErr != nil {
			log.Printf("error !!! %s", execErr.Error())
			_ = tx.Rollback()
			return execErr
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	log.Printf("%s executed!\n", file_name)
	return nil
}

// Функция получает из БД данные о ресторанах
func GetRestaurants(db *sql.DB) ([]structures.Restaurant, error) {
	q := `
		SELECT * FROM restaurants;
	`
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rests []structures.Restaurant
	for rows.Next() {
		var rest structures.Restaurant
		if err := rows.Scan(&rest.ID, &rest.Name, &rest.Avg_time, &rest.Avg_price); err != nil {
			return nil, err
		}
		rests = append(rests, rest)
	}
	return rests, nil
}

// Функция получает из БД столики, которые свободны на
// момент t
func GetTablesByTime(db *sql.DB, t time.Time) ([]structures.Table, error) {
	log.Printf("\n-----* start formattin t (%v)*-----\n", t)
	ts := pq.FormatTimestamp(t)
	log.Printf("\n-----* t formatted %v *-----\n", string(ts))
	q := `
		SELECT DISTINCT tables.id, tables.restaurant_id, tables.capacity 
		FROM tables LEFT OUTER JOIN reservations 
		ON (tables.id = reservations.table_id) WHERE
		($1::timestamp + '2 hours'::interval <= start_time) OR 
		($1::timestamp >= end_time) OR 
		(reservations.id IS NULL) 
		ORDER BY tables.id;
	`

	rows, err := db.Query(q, ts)
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

// Функция добавляет новую бронь в БД
func PostReservations(db *sql.DB, res structures.RawReservation) (structures.Reservation, error) {

	log.Printf("\nOh, hello, let's go, new reservation %v\n", res)

	ctx := context.Background()
	squery := `
		SELECT * FROM clients  WHERE
		(name=$1) AND 
		(phone_number=$2);
	`

	log.Printf("\nDo some query.........\n")

	var client structures.Client
	row := db.QueryRow(squery, res.Reserved_by.Name, res.Reserved_by.Phone)

	log.Printf("\nStart scanning!\n")

	if err := row.Scan(&client.ID, &client.Name, &client.Phone); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("\nOh, no such user!\n")
			client, err = AddClient(db, res.Reserved_by)
			if err != nil {
				log.Printf("\nError when add new user!\n")
				return structures.Reservation{}, err
			}
		} else {
			return structures.Reservation{}, err
		}
	}

	log.Printf("\nGet client struct (%v)\nStart txn\n", client)

	tx, err := db.BeginTx(ctx, nil)
	// defer func() (structures.Reservation, error) {
	// 	if r := recover(); r != nil {
	// 		_ = tx.Rollback()
	// 		return structures.Reservation{}, errors.New("Fail txn!")
	// 	} else {
	// 		return structures.Reservation{}, nil
	// 	}
	// }()

	if err != nil {
		return structures.Reservation{}, err
	}
	q := `
	INSERT INTO reservations (table_id,start_time,end_time,reserved_by)
	VALUES  ($1, $2::timestamp, $2::timestamp + '2 H'::interval, $3);
	`
	t, err := time.Parse(time.RFC3339, res.Start_time)
	_, execErr := tx.Exec(q, res.Table_id, pq.FormatTimestamp(t), client.ID)
	if execErr != nil {
		_ = tx.Rollback()
		return structures.Reservation{}, execErr
	}
	if err := tx.Commit(); err != nil {
		return structures.Reservation{}, err
	}

	var addedReservation structures.Reservation

	row = db.QueryRow(`
		SELECT * FROM reservations
		WHERE (table_id=$1) 
		AND (start_time=$2::timestamp) 
		AND (reserved_by=$3);
	`, res.Table_id, pq.FormatTimestamp(t), client.ID)

	if err := row.Scan(&addedReservation.ID,
		&addedReservation.Table_id,
		&addedReservation.Start_time,
		&addedReservation.End_time,
		&addedReservation.Reserved_by); err != nil {
		return structures.Reservation{}, err
	}

	return addedReservation, nil
}

// Функция добавляет нового клиента в БД
func AddClient(db *sql.DB, rclient structures.RawClient) (structures.Client, error) {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return structures.Client{}, err
	}
	q := `
	INSERT INTO clients (name, phone_number)
	VALUES  ($1, $2);
	`
	_, execErr := tx.Exec(q, rclient.Name, rclient.Phone)
	if execErr != nil {
		_ = tx.Rollback()
		return structures.Client{}, execErr
	}
	if err := tx.Commit(); err != nil {
		return structures.Client{}, err
	}

	squery := `
		SELECT last_value FROM clients_id_seq;
	`
	row := db.QueryRow(squery)
	var client_id int
	if err := row.Scan(&client_id); err != nil {
		return structures.Client{}, err
	}
	return structures.Client{
		ID:    fmt.Sprintf("%v", client_id),
		Name:  rclient.Name,
		Phone: rclient.Phone}, nil
}

// Функция получет из БД данные о клиенте
// по его имени и номеру телефона
func GetClientByRaw(db *sql.DB, rclient structures.RawClient) (structures.Client, error) {
	squery := `
	SELECT * FROM clients  WHERE
	(name=$1) AND 
	(phone_number=$2);
	`
	var client structures.Client
	row := db.QueryRow(squery, rclient.Name, rclient.Phone)
	if err := row.Scan(&client.ID, &client.Name, &client.Phone); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			client, err := AddClient(db, rclient)
			if err != nil {
				return structures.Client{}, err
			}
			return client, nil
		}
		return structures.Client{}, err
	}
	return client, nil
}
