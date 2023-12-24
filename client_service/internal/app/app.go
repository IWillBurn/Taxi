package app

import (
	"context"
	"log"
	"net/http"
	"offering_service/internal/config"
	"offering_service/internal/httpadapter"
	"offering_service/internal/service"
	"os"
	"os/signal"
	"syscall"
)

type app struct {
	config *config.Config

	httpAdapter *httpadapter.Adapter
}

func (a *app) Serve() error {
	done := make(chan os.Signal, 1)

	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := a.httpAdapter.Serve(); err != nil && err != http.ErrServerClosed {
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

	offeringService := &service.LinearOfferingService{LinearCost: 1, BaseCost: 100, PlanetRadius: 6370}
	signingService := &service.JWTSigningService{Key: config.Singing.Key}

	a := &app{
		config:          config,
		offeringService: offeringService,
		signingService:  signingService,
		httpAdapter:     httpadapter.New(&config.HTTP, offeringService, signingService),
	}

	return a, nil
}
