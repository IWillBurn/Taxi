package main

import (
	"client_service/internal/app"
	"client_service/internal/config"
	"flag"
	"fmt"
)

func getConfigPath() string {
	var configPath string

	flag.StringVar(&configPath, "c", "../.config/app.yaml", "path to config file")
	flag.Parse()

	return configPath
}

func main() {

	cfg, err := config.NewConfig(getConfigPath())
	fmt.Println(cfg)
	if err != nil {
	}

	a, err := app.New(cfg)
	if err != nil {
	}

	if err := a.Serve(); err != nil {
	}
}
