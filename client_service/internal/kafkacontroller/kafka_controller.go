package kafkacontroller

import (
	"client_service/internal/metrics"
	"client_service/internal/models"
	"client_service/internal/repo"
	"client_service/internal/socketlistener/publishers"
	"context"
	"encoding/json"
	"log"
	"sync"
)

type KafkaController struct {
	Repo            repo.DataBaseController
	StatusPublisher *publishers.Publisher
	Connection      *Connection
	Waiters         sync.Map

	metrics *metrics.Metrics
}

func (kafkaController *KafkaController) updateStatus(message []byte, status string) {
	var responseMessage models.EventTripStatusUpdate
	err := json.Unmarshal(message, &responseMessage)
	if err != nil {
		log.Fatal(err)
	}
	err = kafkaController.Repo.ChangeTripByTripId(responseMessage.TripId, "status", status)
	if err != nil {
		return
	}
	trip, err := kafkaController.Repo.GetTripByTripId(responseMessage.TripId)
	if err != nil {
		return
	}
	data := make(map[string]string)
	data["status"] = status
	kafkaController.StatusPublisher.Publish(trip.ClientId, data)
}

func (kafkaController *KafkaController) acceptTrip(message []byte, id string) {
	kafkaController.metrics.InTheQueueCounter.Inc()
	kafkaController.updateStatus(message, "DRIVER_FOUND")
}
func (kafkaController *KafkaController) cancelTrip(message []byte, id string) {
	kafkaController.metrics.CanceledTripCounter.Inc()

	var responseMessage models.EventTripStatusUpdate
	err := json.Unmarshal(message, &responseMessage)
	if err != nil {
		log.Fatal(err)
	}
	trip, err := kafkaController.Repo.GetTripByTripId(responseMessage.TripId)
	if err != nil {
		log.Fatal(err)
	}
	if trip.Status == "DRIVER_FOUND" {
		kafkaController.metrics.InTheQueueCounter.Dec()
	}
	kafkaController.updateStatus(message, "CANCELED")
}
func (kafkaController *KafkaController) createTrip(message []byte, id string) {
	log.Println("GOT IT!")
	var responseMessage models.EventCreateTrip
	err := json.Unmarshal(message, &responseMessage)
	if err != nil {
		log.Fatal(err)
	}
	err = kafkaController.Repo.ChangeTripByOfferId(responseMessage.OfferId, "status", responseMessage.Status)
	if err != nil {
		return
	}
	err = kafkaController.Repo.ChangeTripByOfferId(responseMessage.OfferId, "trip_id", responseMessage.TripId)
	if err != nil {
		return
	}

	trip, err := kafkaController.Repo.GetTripByTripId(responseMessage.TripId)
	if err != nil {
		return
	}
	data := make(map[string]string)
	data["status"] = responseMessage.Status

	kafkaController.metrics.CreatedOrdersCounter.Inc()
	log.Println(responseMessage.TripId)
	kafkaController.StatusPublisher.Publish(trip.ClientId, data)
	_, ok := kafkaController.Waiters.Load(id)
	if ok {
		kafkaController.Waiters.Delete(id)
		kafkaController.Waiters.Store(id, responseMessage.TripId)
	}
}
func (kafkaController *KafkaController) startTrip(message []byte, id string) {
	kafkaController.metrics.InTheQueueCounter.Dec()
	kafkaController.updateStatus(message, "STARTED")
}
func (kafkaController *KafkaController) endTrip(message []byte, id string) {
	kafkaController.metrics.EndedTripCounter.Inc()
	kafkaController.updateStatus(message, "ENDED")
}

func NewService(connection *Connection, metrics *metrics.Metrics) *KafkaController {

	k := &KafkaController{
		Connection: connection,
		metrics:    metrics,
	}

	k.metrics.Serve()
	connection.AddHandler("trip.event.created", k.createTrip)
	connection.AddHandler("trip.event.accepted", k.acceptTrip)
	connection.AddHandler("trip.event.canceled", k.cancelTrip)
	connection.AddHandler("trip.event.ended", k.endTrip)
	connection.AddHandler("trip.event.started", k.startTrip)
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
