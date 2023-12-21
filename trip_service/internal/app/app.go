package app

import (
	"context"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"trip_service/internal/app/config"
	"trip_service/internal/app/service"
	"trip_service/internal/connection"
)

type App struct {
	config      *config.Config
	connection  *connection.Connection
	tripService *service.TripService
}

func NewApp(config *config.Config) *App {
	return &App{
		config: config,
		connection: connection.NewConnection(
			&kafka.ReaderConfig{
				Topic:   config.Connection.Inbound.Topic,
				Brokers: config.Connection.Inbound.Brokers,
			}, &kafka.WriterConfig{
				Topic:   config.Connection.Outbound.Topic,
				Brokers: config.Connection.Outbound.Brokers,
			}),
	}
}
func (app *App) Run(ctx context.Context) error {
	done := make(chan os.Signal, 1)

	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := app.connection.Serve(ctx); err != nil {
			log.Fatal("Error on start server", zap.Error(err))
		}
	}()

	<-done

	err := app.Shutdown()
	if err != nil {
		log.Fatal("Error on shutdown service", err)
		return err
	}
	return nil
}

func (app *App) Shutdown() error {
	_, cancel := context.WithTimeout(context.Background(), app.config.App.ShutdownTimeout)
	defer cancel()

	err := app.connection.Close()
	return err
}
