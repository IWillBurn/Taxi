package connection

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"trip_service/internal/app/service"
	"trip_service/internal/model"
)

type Connection struct {
	reader *kafka.Reader
	writer *kafka.Writer

	tripService *service.TripService
	handlers    map[string]func([]byte)
}

func (connection *Connection) addHandler(contentType string, reducer func([]byte)) {
	connection.handlers[contentType] = reducer
}
func (connection *Connection) handle(contentType string, message []byte) {
	connection.handlers[contentType](message)
}
func (connection *Connection) acceptTrip(message []byte) {
	var responseMessage model.AcceptTrip
	err := json.Unmarshal(message, &responseMessage)
	if err != nil {
		log.Fatal(err)
	}
	connection.tripService.AcceptTrip(responseMessage)
}
func (connection *Connection) cancelTrip(message []byte) {
	var responseMessage model.CancelTrip
	err := json.Unmarshal(message, &responseMessage)
	if err != nil {
		log.Fatal(err)
	}
	connection.tripService.CancelTrip(responseMessage)
}
func (connection *Connection) creatTrip(message []byte) {
	var responseMessage model.CreatTrip
	err := json.Unmarshal(message, &responseMessage)
	if err != nil {
		log.Fatal(err)
	}
	connection.tripService.CreatTrip(responseMessage)
}
func (connection *Connection) endTrip(message []byte) {
	var responseMessage model.EndTrip
	err := json.Unmarshal(message, &responseMessage)
	if err != nil {
		log.Fatal(err)
	}
	connection.tripService.EndTrip(responseMessage)
}
func (connection *Connection) startTrip(message []byte) {
	var responseMessage model.StartTrip
	err := json.Unmarshal(message, &responseMessage)
	if err != nil {
		log.Fatal(err)
	}
	connection.tripService.StartTrip(responseMessage)
}

func NewConnection(readerConfig *kafka.ReaderConfig, writerConfig *kafka.WriterConfig) *Connection {
	connection := Connection{
		reader:   kafka.NewReader(*readerConfig),
		writer:   kafka.NewWriter(*writerConfig),
		handlers: make(map[string]func([]byte)),
	}

	connection.addHandler("trip.command.accept", connection.acceptTrip)
	connection.addHandler("trip.command.cancel", connection.cancelTrip)
	connection.addHandler("trip.command.creat", connection.creatTrip)
	connection.addHandler("trip.command.end", connection.endTrip)
	connection.addHandler("trip.command.start", connection.startTrip)

	return &connection
}

func (connection *Connection) Serve(ctx context.Context) error {
	log.Print("start reading")
	for {
		msg, err := connection.Read(ctx)
		if err != nil {
			log.Println(err)
		}
		var responseMessage model.InboundMessage
		err = json.Unmarshal(msg, &responseMessage)
		if err != nil {
			log.Println(err)
		}
		data, err := json.Marshal(responseMessage.Data)
		if err != nil {
			log.Println(err)
		}
		connection.handle(responseMessage.DataContentType, data)
	}
}

func (connection *Connection) Write(ctx context.Context, key []byte, msg []byte) error {
	err := connection.writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: msg,
	})
	return err
}
func (connection *Connection) Read(ctx context.Context) ([]byte, error) {
	msg, err := connection.reader.ReadMessage(ctx)
	if err != nil {
		log.Print("Problem with reading message", err)
		return []byte{}, err
	}
	return msg.Value, err
}
func (connection *Connection) Close() error {
	errWriter := connection.writer.Close()
	errReader := connection.reader.Close()

	if errWriter != nil {
		return errWriter
	}
	if errReader != nil {
		return errReader
	}
	return nil
}
