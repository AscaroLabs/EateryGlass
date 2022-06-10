package drop

// Данный пакет отвечает за очищение БД
// при завершении работы приложения

import (
	"log"

	"github.com/AscaroLabs/EateryGlass/internal/config"
	"github.com/AscaroLabs/EateryGlass/pkg/database/postgresql/storage"
	_ "github.com/lib/pq"
)

// Функция, очищающая БД
// Вызывается при получении сигнала SIGINT
func Cleanup(cfg *config.Config) func() {
	db := storage.NewDB(cfg)
	return func() {
		err := storage.ExecFromFile(cfg, db, "drop_tables.sql")
		if err != nil {
			log.Fatal(err)
		}
	}
}
