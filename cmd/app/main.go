package main

import (
	"github.com/AscaroLabs/EateryGlass/internal/app"
	"github.com/AscaroLabs/EateryGlass/internal/config"
	"github.com/AscaroLabs/EateryGlass/pkg/database/postgresql/drop"
	"github.com/xlab/closer"
)

func main() {

	cfg := config.NewConfig()
	closer.Bind(drop.Cleanup(cfg))
	// defer storage.ExecFromFile(cfg, db, "drop_tables.sql")
	go app.Run(cfg)
	closer.Hold()
}
