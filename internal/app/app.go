package app

// Данный пакет отвечает за запуск сервера
// и инициализации БД

import (
	"fmt"

	"github.com/AscaroLabs/EateryGlass/internal/config"
	"github.com/AscaroLabs/EateryGlass/internal/delivery"
	"github.com/AscaroLabs/EateryGlass/pkg/database/postgresql/storage"
)

// Run запускает приложение
func Run(cfg *config.Config) {

	// Инициализируем новый экземпляр sql.BD
	db := storage.NewDB(cfg)

	// Инициализируем новый экземпляр gin.Engine
	router := delivery.NewRouter(db)

	// Запускаем сервер
	router.Run(fmt.Sprintf("%s:%s", cfg.App_host, cfg.App_port))
}
