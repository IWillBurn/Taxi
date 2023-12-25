package service

import (
	"client_service/internal/httpadapter/responses"
	"client_service/internal/kafkacontroller"
	"client_service/internal/models"
	"client_service/internal/repo"
	"client_service/internal/socketlistener/publishers"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type DefaultTripService struct {
	KafkaController     *kafkacontroller.KafkaController
	DataBaseController  repo.DataBaseController
	OfferingServiceHost string
	client              http.Client
}

func (tripService *DefaultTripService) CreateTrip(offerId string) error {
	id := uuid.New().String()
	err := tripService.KafkaController.Connection.Write(context.Background(), []byte(id),
		models.OutboundMessage{
			Id:              id,
			Source:          "/client",
			Type:            "trip.command.create",
			DataContentType: "application/json",
			Time:            time.Now().Format(time.RFC3339),
			Data: models.CreateData{
				OfferId: offerId,
			},
		})
	if err != nil {
		log.Fatal(err)
	}

	resp, err := tripService.client.Get(tripService.OfferingServiceHost + "/offers/" + offerId)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var response responses.OfferResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	trip := repo.Trip{
		InitId:   id,
		TripId:   "",
		OfferId:  offerId,
		From:     response.From,
		To:       response.To,
		ClientId: response.ClientId,
		Price:    response.Price,
		Status:   "",
	}

	err = (tripService.DataBaseController).AddTrip(trip)
	if err != nil {
		return err
	}
	return nil
}

func (tripService *DefaultTripService) CancelTrip(tripId string, reason string) error {
	id := uuid.New().String()
	err := tripService.KafkaController.Connection.Write(context.Background(), []byte(id),
		models.OutboundMessage{
			Id:              id,
			Source:          "/client",
			Type:            "trip.command.cancel",
			DataContentType: "application/json",
			Time:            time.Now().Format(time.RFC3339),
			Data: models.CancelData{
				TripId: tripId,
				Reason: reason,
			},
		})
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (tripService *DefaultTripService) GetTripStatus(clientId string, tripId string, publisher *publishers.Publisher) error {
	fmt.Println("To publish")
	publisher.Publish(clientId, "OK")
	return nil
}