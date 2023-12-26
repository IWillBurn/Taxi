package app

import (
	"client_service/internal/config"
	"client_service/internal/httpadapter"
	"client_service/internal/kafkacontroller"
	"client_service/internal/metrics"
	"client_service/internal/repo"
	"client_service/internal/service"
	"client_service/internal/socketlistener"
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type app struct {
	config             *config.Config
	dataBaseController repo.DataBaseController
	tripService        service.TripService
	socketController   *socketlistener.SocketController
	httpAdapter        *httpadapter.Adapter
	kafkaController    *kafkacontroller.KafkaController
}

func (a *app) Serve() error {
	done := make(chan os.Signal, 1)

	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := a.httpAdapter.Serve(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err.Error())
		}
	}()

	go func() {
		if err := a.socketController.Serve(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err.Error())
		}
	}()

	<-done

	a.Shutdown()

	return nil
}

func (a *app) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), a.config.App.ShutdownTimeout)
	defer cancel()

	a.httpAdapter.Shutdown(ctx)
}

func New(config *config.Config) (App, error) {

	DataBaseController := &repo.MongoDB{Config: &config.Mongo}
	err := DataBaseController.Serve()
	if err != nil {
		return nil, err
	}

	connection := kafkacontroller.NewConnection(
		&kafka.ReaderConfig{
			Topic:   config.Connection.Outbound.Topic,
			Brokers: config.Connection.Outbound.Brokers,
		},
		&kafka.WriterConfig{
			Topic:   config.Connection.Inbound.Topic,
			Brokers: config.Connection.Inbound.Brokers,
		})

	kafkaController := kafkacontroller.NewService(connection, metrics.NewMetrics())
	go kafkaController.Serve(context.Background())

	tripService := &service.DefaultTripService{
		KafkaController:     kafkaController,
		OfferingServiceHost: config.HTTP.OfferingAddress,
		DataBaseController:  DataBaseController,
		Client:              http.Client{},
	}

	socketController, _ := socketlistener.NewSocketController(&config.Socket, tripService)
	kafkaController.Repo = DataBaseController
	kafkaController.StatusPublisher = (*socketController.Publishers)["status"]
	go kafkaController.Serve(context.Background())
	a := &app{
		config:             config,
		dataBaseController: DataBaseController,
		tripService:        tripService,
		socketController:   socketController,
		httpAdapter:        httpadapter.New(&config.HTTP, DataBaseController, tripService),
		kafkaController:    kafkaController,
	}

	return a, nil
}
