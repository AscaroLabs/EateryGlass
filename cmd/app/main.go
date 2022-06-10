package main

import (
	"github.com/AscaroLabs/EateryGlass/internal/app"
	"github.com/AscaroLabs/EateryGlass/internal/config"
	"github.com/AscaroLabs/EateryGlass/pkg/database/postgresql/drop"
	"github.com/xlab/closer"
)

func main() {

	// Создаем новый объект config.Config
	// (т.е. читаем переменные из .env файла)
	cfg := config.NewConfig()
	// Связываем новую функцию с closer'ом
	// Это позволит очищать БД при завершении работы
	// контейнера при получении сигнала SIGINT (ctrl+C)
	closer.Bind(drop.Cleanup(cfg))
	// Запускаем приложение в отдельной горутине, что бы
	// Основная горутина могла обрабатывать входящие сигналы
	go app.Run(cfg)
	// Блокируем основную горутину на получение сигналов
	closer.Hold()
}
