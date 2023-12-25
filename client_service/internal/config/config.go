package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

const (
	AppName                   = "offering"
	DefaultServeAddress       = "localhost:32434"
	DefaultShutdownTimeout    = 20 * time.Second
	DefaultBasePath           = "/"
	DefaultSocketPath         = "/"
	DefaultSocketServeAddress = "localhost:53242"
	DefaultURI                = "mongodb://localhost:27017"
)

type AppConfig struct {
	Debug           bool          `yaml:"debug"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type HttpAdapterConfig struct {
	ServeAddress    string `yaml:"serve_address"`
	BasePath        string `yaml:"base_path"`
	OfferingAddress string `yaml:"offering_address"`
}

type SocketConfig struct {
	ServeAddress string `yaml:"socket_serve_address"`
	BasePath     string `yaml:"socket_path"`
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

type MongoConfig struct {
	URI string `yaml:"uri"`
}

type Config struct {
	App        AppConfig         `yaml:"app"`
	HTTP       HttpAdapterConfig `yaml:"http"`
	Socket     SocketConfig      `yaml:"socket"`
	Mongo      MongoConfig       `yaml:"mongo"`
	Connection ConnectionConfig  `yaml:"connection"`
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
		Socket: SocketConfig{
			ServeAddress: DefaultSocketServeAddress,
			BasePath:     DefaultSocketPath,
		},
		Mongo: MongoConfig{
			URI: DefaultURI,
		},
		Connection: ConnectionConfig{
			Inbound:  InboundConfig{},
			Outbound: OutboundConfig{},
		},
	}

	if err := yaml.Unmarshal(data, &cnf); err != nil {
		return nil, err
	}

	return &cnf, nil
}

/*
func DefaultConfig() (*Config, error) {
	cnf := Config{
		App: AppConfig{
			ShutdownTimeout: DefaultShutdownTimeout,
		},
		HTTP: HttpAdapterConfig{
			ServeAddress: DefaultServeAddress,
			BasePath:     DefaultBasePath,
		},
		Socket: SocketConfig{
			ServeAddress: DefaultSocketServeAddress,
			BasePath:     DefaultSocketPath,
		},
		Mongo: MongoConfig{
			URI: DefaultURI,
		},
	}
	return &cnf, nil
}
*/
