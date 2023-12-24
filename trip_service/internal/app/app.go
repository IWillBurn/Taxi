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
	"trip_service/internal/repository"
)

type App struct {
	config      *config.Config
	connection  *connection.Connection
	tripService *service.TripService
}

func NewApp(ctx context.Context, config *config.Config) *App {
	//Create topics by docker-compose
	//topic.CreateTopic(ctx, config.Connection.Inbound.Brokers, kafka.TopicConfig{
	//	Topic: config.Connection.Inbound.Topic,
	//})
	//topic.CreateTopic(ctx, config.Connection.Outbound.Brokers, kafka.TopicConfig{
	//	Topic: config.Connection.Outbound.Topic,
	//})
	log.Println("topic created")
	appConnection := connection.NewConnection(
		&kafka.ReaderConfig{
			Topic:   config.Connection.Inbound.Topic,
			Brokers: config.Connection.Inbound.Brokers,
		}, &kafka.WriterConfig{
			Topic:   config.Connection.Outbound.Topic,
			Brokers: config.Connection.Outbound.Brokers,
		})
	return &App{
		config:      config,
		connection:  appConnection,
		tripService: service.NewService(repository.NewRepository(), appConnection),
	}
}
func (app *App) Run(ctx context.Context) error {
	done := make(chan os.Signal, 1)

	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := app.tripService.Serve(ctx); err != nil {
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

	err := app.tripService.Shutdown()
	return err
}
