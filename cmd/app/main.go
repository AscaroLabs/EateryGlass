package main

import (
	"github.com/AscaroLabs/EateryGlass/internal/app"
	"github.com/AscaroLabs/EateryGlass/internal/config"
)

func main() {
	cfg := config.NewConfig()
	app.Run(cfg)
}
