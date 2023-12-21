package main

import (
	"context"
	"flag"
	"go.uber.org/zap"
	"log"
	"trip_service/internal/app"
	"trip_service/internal/app/config"
)

func getConfigPath() string {
	var configPath string

	flag.StringVar(&configPath, "c", "./.config/trip_service.yaml", "path to config file")
	flag.Parse()

	return configPath
}

func main() {
	cfg, err := config.NewConfig(getConfigPath())
	if err != nil {
		log.Fatalln("Error on read config", err)
	}
	application := app.NewApp(cfg)

	if err := application.Run(context.Background()); err != nil {
		log.Fatal("Error on server", zap.Error(err))
	}
}
