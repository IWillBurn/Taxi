package config

import (
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

type DatabaseConfig struct {
	DSN           string `yaml:"dsn"`
	MigrationsDir string `yaml:"migrations_dir"`
}
type InboundConfig struct {
	Topic   string   `yaml:"topic"`
	Brokers []string `yaml:"brokers"`
}
type OutboundConfig struct {
	Topic   string   `yaml:"topic"`
	Brokers []string `yaml:"brokers"`
}
type ConnectionConfig struct {
	Inbound  InboundConfig  `yaml:"inbound"`
	Outbound OutboundConfig `yaml:"outbound"`
}
type AppConfig struct {
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}
type Config struct {
	Connection ConnectionConfig `yaml:"connection"`
	Database   DatabaseConfig   `yaml:"database"`
	App        AppConfig        `yaml:"app"`
}

func NewConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
