package kafkacontroller

import (
	"client_service/internal/models"
	"client_service/internal/repo"
	"client_service/internal/socketlistener/publishers"
	"context"
	"encoding/json"
	"log"
)

type KafkaController struct {
	repo            repo.DataBaseController
	statusPublisher *publishers.Publisher
	Connection      *Connection
}

func (kafkaController *KafkaController) updateStatus(message []byte, status string) {
	var responseMessage models.EventTripStatusUpdate
	err := json.Unmarshal(message, &responseMessage)
	if err != nil {
		log.Fatal(err)
	}
	err = kafkaController.repo.ChangeTripByOfferId(responseMessage.TripId, "status", status)
	if err != nil {
		return
	}
	trip, err := kafkaController.repo.GetTripByTripId(responseMessage.TripId)
	if err != nil {
		return
	}
	data := make(map[string]string)
	data["status"] = status
	kafkaController.statusPublisher.Publish(trip.ClientId, data)
}

func (kafkaController *KafkaController) acceptTrip(message []byte) {
	kafkaController.updateStatus(message, "DRIVER_FOUND")
}
func (kafkaController *KafkaController) cancelTrip(message []byte) {
	kafkaController.updateStatus(message, "CANCELED")
}
func (kafkaController *KafkaController) createTrip(message []byte) {
	var responseMessage models.EventCreateTrip
	err := json.Unmarshal(message, &responseMessage)
	if err != nil {
		log.Fatal(err)
	}
	err = kafkaController.repo.ChangeTripByOfferId(responseMessage.OfferId, "status", responseMessage.Status)
	if err != nil {
		return
	}
	err = kafkaController.repo.ChangeTripByOfferId(responseMessage.OfferId, "trip_id", responseMessage.TripId)
	if err != nil {
		return
	}

	trip, err := kafkaController.repo.GetTripByTripId(responseMessage.TripId)
	if err != nil {
		return
	}
	data := make(map[string]string)
	data["status"] = responseMessage.Status
	kafkaController.statusPublisher.Publish(trip.ClientId, data)
}
func (kafkaController *KafkaController) startTrip(message []byte) {
	kafkaController.updateStatus(message, "STARTED")
}
func (kafkaController *KafkaController) endTrip(message []byte) {
	kafkaController.updateStatus(message, "ENDED")
}

func NewService(connection *Connection) *KafkaController {

	k := &KafkaController{
		Connection: connection,
	}
	connection.AddHandler("trip.event.accept", k.acceptTrip)
	connection.AddHandler("trip.event.cancel", k.cancelTrip)
	connection.AddHandler("trip.event.end", k.endTrip)
	connection.AddHandler("trip.event.start", k.startTrip)
	return k
}
func (kafkaController *KafkaController) Serve(ctx context.Context) error {
	err := kafkaController.Connection.Serve(ctx)
	if err != nil {
		return err
	}
	return nil
}
func (kafkaController *KafkaController) Shutdown() error {
	return kafkaController.Connection.Close()
}
