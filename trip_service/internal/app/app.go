package app

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
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

const driverName = "postgres"

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

	db, err := initDB(ctx, &config.Database)
	if err != nil {
		log.Fatal("Problem with init database ", err)
	}
	appRepo := repository.NewRepository(db)

	return &App{
		config:      config,
		connection:  appConnection,
		tripService: service.NewService(appRepo, appConnection),
	}
}

func initDB(ctx context.Context, config *config.DatabaseConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open(driverName, config.DSN)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	// migrations
	fs := os.DirFS(config.MigrationsDir)
	goose.SetBaseFS(fs)

	if err = goose.SetDialect(driverName); err != nil {
		panic(err)
	}

	if err = goose.UpContext(ctx, db.DB, "."); err != nil {
		panic(err)
	}

	return db, nil
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
