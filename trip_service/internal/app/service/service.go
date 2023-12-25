package service

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"trip_service/internal/connection"
	"trip_service/internal/model"
)

type TripService struct {
	repo       model.TripRepository
	connection *connection.Connection
}

func (service *TripService) acceptTrip(message []byte) {
	var responseMessage model.AcceptTrip
	err := json.Unmarshal(message, &responseMessage)
	if err != nil {
		log.Fatal(err)
	}
	err = service.repo.UpdateStarted(
		context.Background(),
		&model.ParamsStarted{
			Id: responseMessage.TripID,
		},
		&model.TripStarted{
			DriverId:     responseMessage.DriverID,
			CurrentStage: "DRIVER FOUND",
		})
	if err != nil {
		log.Fatal(err)
	}
	err = service.connection.Write(context.Background(), []byte(responseMessage.TripID+"DRIVER FOUND"), model.OutboundMessage{
		ID:              responseMessage.TripID,
		Source:          "/trip",
		Type:            "trip.event.accepted",
		DataContentType: "application/json",
		Time:            time.Now().Format(time.RFC3339),
		Data: model.EventAcceptTrip{
			TripID:   responseMessage.TripID,
			DriverID: responseMessage.DriverID,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
func (service *TripService) cancelTrip(message []byte) {
	var responseMessage model.CancelTrip
	err := json.Unmarshal(message, &responseMessage)
	if err != nil {
		log.Fatal(err)
	}
	trips, err := service.repo.GetStarted(context.Background(),
		&model.ParamsStarted{
			Id: responseMessage.TripID,
		})
	if (err != nil) || (len(trips) != 1) {
		log.Fatal(err)
	}
	_, err = service.repo.CreateFinished(context.Background(),
		&model.TripFinished{
			TripStarted: trips[0],
			Successful:  false,
			Reason:      responseMessage.Reason,
		})
	err = service.connection.Write(context.Background(), []byte(responseMessage.TripID+"DRIVER FOUND"), model.OutboundMessage{
		ID:              responseMessage.TripID,
		Source:          "/trip",
		Type:            "trip.event.canceled",
		DataContentType: "application/json",
		Time:            time.Now().Format(time.RFC3339),
		Data: model.EventCancelTrip{
			TripID: responseMessage.TripID,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
func (service *TripService) createTrip(message []byte) {
	var responseMessage model.CreatTrip
	err := json.Unmarshal(message, &responseMessage)
	if err != nil {
		log.Fatal(err)
	}
	tripId, err := service.repo.CreateStarted(context.Background(),
		&model.TripStarted{
			OfferId:      responseMessage.OfferID,
			CurrentStage: "CREATED",
		})
	resp, err := http.Get("https://httpbin.org/get")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
	err = service.connection.Write(context.Background(), []byte(tripId+"CREATED"), model.OutboundMessage{
		ID:              tripId,
		Source:          "/trip",
		Type:            "trip.event.canceled",
		DataContentType: "application/json",
		Time:            time.Now().Format(time.RFC3339),
		Data: model.EventCreatTrip{
			OfferID: responseMessage.OfferID,
		},
	})
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
	err = service.repo.UpdateStarted(
		context.Background(),
		&model.ParamsStarted{
			Id: responseMessage.TripID,
		},
		&model.TripStarted{
			CurrentStage: "STARTED",
		})
	err = service.connection.Write(context.Background(), []byte(responseMessage.TripID+"STARTED"), model.OutboundMessage{
		ID:              responseMessage.TripID,
		Source:          "/trip",
		Type:            "trip.event.started",
		DataContentType: "application/json",
		Time:            time.Now().Format(time.RFC3339),
		Data: model.EventStartTrip{
			TripID: responseMessage.TripID,
		},
	})
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
	trips, err := service.repo.GetStarted(context.Background(),
		&model.ParamsStarted{
			Id: responseMessage.TripID,
		})
	if (err != nil) || (len(trips) != 1) {
		log.Fatal(err)
	}
	_, err = service.repo.CreateFinished(context.Background(),
		&model.TripFinished{
			TripStarted: trips[0],
			Successful:  true,
		})
	err = service.repo.DeleteStarted(context.Background(), &model.ParamsStarted{Id: responseMessage.TripID})
	err = service.connection.Write(context.Background(), []byte(responseMessage.TripID+"ENDED"), model.OutboundMessage{
		ID:              responseMessage.TripID,
		Source:          "/trip",
		Type:            "trip.event.ended",
		DataContentType: "application/json",
		Time:            time.Now().Format(time.RFC3339),
		Data: model.EventEndTrip{
			TripID: responseMessage.TripID,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

func NewService(
	repo model.TripRepository,
	connection *connection.Connection) *TripService {

	service := &TripService{
		repo:       repo,
		connection: connection,
	}
	connection.AddHandler("trip.command.accept", service.acceptTrip)
	connection.AddHandler("trip.command.cancel", service.cancelTrip)
	connection.AddHandler("trip.command.create", service.createTrip)
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
