package service

import (
	"context"
	"encoding/json"
	"log"
	"trip_service/internal/connection"
	"trip_service/internal/model"
)

type TripService struct {
	repo       *model.TripRepository
	connection *connection.Connection
}

func (service *TripService) acceptTrip(message []byte) {
	var responseMessage model.AcceptTrip
	err := json.Unmarshal(message, &responseMessage)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Accept Trip")

}
func (service *TripService) cancelTrip(message []byte) {
	var responseMessage model.CancelTrip
	err := json.Unmarshal(message, &responseMessage)
	if err != nil {
		log.Fatal(err)
	}
}
func (service *TripService) creatTrip(message []byte) {
	var responseMessage model.CreatTrip
	err := json.Unmarshal(message, &responseMessage)
	if err != nil {
		log.Fatal(err)
	}
}
func (service *TripService) endTrip(message []byte) {
	var responseMessage model.EndTrip
	err := json.Unmarshal(message, &responseMessage)
	if err != nil {
		log.Fatal(err)
	}
}
func (service *TripService) startTrip(message []byte) {
	var responseMessage model.StartTrip
	err := json.Unmarshal(message, &responseMessage)
	if err != nil {
		log.Fatal(err)
	}

}

func NewService(
	repo model.TripRepository,
	connection *connection.Connection) *TripService {

	service := &TripService{
		repo:       &repo,
		connection: connection,
	}
	connection.AddHandler("trip.command.accept", service.acceptTrip)
	connection.AddHandler("trip.command.cancel", service.cancelTrip)
	connection.AddHandler("trip.command.creat", service.creatTrip)
	connection.AddHandler("trip.command.end", service.endTrip)
	connection.AddHandler("trip.command.start", service.startTrip)
	return service
}
func (service *TripService) Serve(ctx context.Context) error {
	err := service.connection.Serve(ctx)
	if err != nil {
		return err
	}
	return nil
}
func (service *TripService) Shutdown() error {
	return service.connection.Close()
}
