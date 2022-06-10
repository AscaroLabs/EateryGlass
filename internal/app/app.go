package app

// Данный пакет отвечает за запуск сервера
// и инициализации БД

import (
	"fmt"
	"log"

	"github.com/AscaroLabs/EateryGlass/internal/config"
	"github.com/AscaroLabs/EateryGlass/internal/delivery"
	"github.com/AscaroLabs/EateryGlass/pkg/database/postgresql/storage"
)

func Run(cfg *config.Config) {
	log.Println("(===============*- CREATE DB -*===============)")
	db := storage.NewDB(cfg)
	log.Println("(===============*- DB CREATED -*===============)")
	log.Println("(===============*- CREATE ROUTER -*===============)")
	router := delivery.NewRouter(db)
	log.Println("(===============*- ROUTER CREATED -*===============)")
	router.Run(fmt.Sprintf("%s:%s", cfg.App_host, cfg.App_port))
}
