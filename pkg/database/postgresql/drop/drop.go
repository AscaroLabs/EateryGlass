package drop

import (
	"context"
	"database/sql"

	// "errors"
	"fmt"
	"log"
	"os"
	"strings"

	// "time"
	"github.com/AscaroLabs/EateryGlass/internal/config"
	"github.com/AscaroLabs/EateryGlass/pkg/database/postgresql/storage"
	_ "github.com/lib/pq"
)

func Cleanup(cfg *config.Config) func() {
	db := storage.NewDB(cfg)
	// err := ExecFromFile(db, "create_model.sql")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	return func() {
		err := storage.ExecFromFile(cfg, db, "drop_tables.sql")
		if err != nil {
			log.Fatal(err)
		}
	}

}

func NewDB() *sql.DB {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", "postgres",
		"qwerty",
		"172.17.0.2",
		"5432",
		"postgres")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	// Создаем таблички

	// Заполняем таблички
	// ExecFromFile(cfg, db, "fill_db.sql")
	return db
}

func ExecFromFile(cfg *config.Config, db *sql.DB, file_name string) error {
	// var ctx context.Context

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

	log.Printf("%v byte readed \n", n)

	queries := strings.Split((string(data[:n])), "\n")

	for i, q := range queries {
		log.Printf("Command %v --> %s \n", i, q)
	}

	tx, err := db.BeginTx(context.Background(), nil)

	if err != nil {
		return err
	}
	for _, q := range queries {
		if q == "" {
			continue
		}
		_, execErr := tx.Exec(q)
		if execErr != nil {
			_ = tx.Rollback()
			return execErr
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
