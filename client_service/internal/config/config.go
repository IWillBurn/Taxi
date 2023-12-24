package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

const (
	AppName                = "offering"
	DefaultServeAddress    = "localhost:63342"
	DefaultShutdownTimeout = 20 * time.Second
	DefaultBasePath        = "/"
)

type AppConfig struct {
	Debug           bool          `yaml:"debug"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type HttpAdapterConfig struct {
	ServeAddress string `yaml:"serve_address"`
	BasePath     string `yaml:"base_path"`
}

type MongoConfig struct {
	Key string `yaml:"secret_key"`
}

type Config struct {
	App   AppConfig         `yaml:"app"`
	HTTP  HttpAdapterConfig `yaml:"http"`
	Mongo MongoConfig       `yaml:"mongo"`
}

func NewConfig(fileName string) (*Config, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	cnf := Config{
		App: AppConfig{
			ShutdownTimeout: DefaultShutdownTimeout,
		},
		HTTP: HttpAdapterConfig{
			ServeAddress: DefaultServeAddress,
			BasePath:     DefaultBasePath,
		},
	}

	if err := yaml.Unmarshal(data, &cnf); err != nil {
		return nil, err
	}

	return &cnf, nil
}
