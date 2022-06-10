package storage

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

// Коннектимся к БД
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

	// Создаем и Заполняем таблички
	FillDB(cfg, db)

	return db
}

func FillDB(cfg *config.Config, db *sql.DB) {
	ExecFromFile(cfg, db, "create_model.sql")
	ExecFromFile(cfg, db, "fill_db.sql")
	ExecFromFile(cfg, db, "insert_clients.sql")
	ExecFromFile(cfg, db, "insert_reservations.sql")
}

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

func PostReservations(db *sql.DB, res structures.RawReservation) (structures.Reservation, error) {
	ctx := context.Background()
	squery := `
		SELECT * FROM clients  WHERE
		(name=?) AND 
		(phone_number=?);
	`
	var client structures.Client
	row := db.QueryRow(squery, res.Reserved_by.Name, res.Reserved_by.Phone)
	if err := row.Scan(&client.ID, &client.Name, &client.Phone); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			client, err = AddClient(db, res.Reserved_by)
			if err != nil {
				return structures.Reservation{}, err
			}
		}
		return structures.Reservation{}, err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return structures.Reservation{}, err
	}
	q := `
	INSERT INTO reservations (table_id,start_time,end_time,reserved_by)
	VALUES  (?, ?, ? + '2 H'::interval, ?);
	`
	t, err := time.Parse(time.RFC3339, res.Start_time)
	_, execErr := tx.Exec(q, res.Table_id, pq.FormatTimestamp(t), pq.FormatTimestamp(t), client.ID)
	if execErr != nil {
		_ = tx.Rollback()
		return structures.Reservation{}, execErr
	}
	if err := tx.Commit(); err != nil {
		return structures.Reservation{}, err
	}

	return structures.Reservation{}, nil
}

func AddClient(db *sql.DB, rclient structures.RawClient) (structures.Client, error) {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return structures.Client{}, err
	}
	q := `
	INSERT INTO clients (name, phone_number)
	VALUES  (?, ?);
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

func GetClientByRaw(db *sql.DB, rclient structures.RawClient) (structures.Client, error) {
	squery := `
	SELECT * FROM clients  WHERE
	(name=?) AND 
	(phone_number=?);
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
