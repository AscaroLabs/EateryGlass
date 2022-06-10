package config

// Данный пакет отвечает за чтения переменных
// из .env файла и создания структуры Config,
// которую можно передовать по указателю в другие
// функции

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_host     string
	DB_port     string
	DB_name     string
	DB_user     string
	DB_password string
	App_host    string
	App_port    string
	Main_dir    string
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return &Config{
		DB_host:     os.Getenv("DB_HOST_ADDR"),
		DB_port:     os.Getenv("DB_HOST_PORT"),
		DB_name:     os.Getenv("DB_NAME"),
		DB_user:     os.Getenv("DB_USERNAME"),
		DB_password: os.Getenv("DB_PASSWORD"),
		App_host:    os.Getenv("APP_HOST"),
		App_port:    os.Getenv("APP_PORT"),
		Main_dir:    os.Getenv("MAIN_DIR"),
	}
}
