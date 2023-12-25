package connection

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"strings"
	"trip_service/internal/model"
)

type Connection struct {
	reader   *kafka.Reader
	writer   *kafka.Writer
	handlers map[string]func([]byte)
}

func (connection *Connection) AddHandler(contentType string, reducer func([]byte)) {
	connection.handlers[contentType] = reducer
}
func (connection *Connection) handle(Type string, message []byte) {
	function, exist := connection.handlers[Type]
	log.Println(Type)
	if exist {
		function(message)
	} else {
		log.Println("Nothing to handle")
	}
}

func NewConnection(readerConfig *kafka.ReaderConfig, writerConfig *kafka.WriterConfig) *Connection {

	connection := Connection{
		reader:   kafka.NewReader(*readerConfig),
		writer:   kafka.NewWriter(*writerConfig),
		handlers: make(map[string]func([]byte)),
	}

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
		err = json.NewDecoder(strings.NewReader(string(msg))).Decode(&responseMessage)

		if err != nil {
			log.Println(err)
		}
		data, err := json.Marshal(responseMessage.Data)
		if err != nil {
			log.Println(err)
		}
		connection.handle(responseMessage.Type, data)
	}
}

func (connection *Connection) Write(ctx context.Context, key []byte, msg model.OutboundMessage) error {
	encodeMsg, err := json.Marshal(msg)
	if err != nil {
		log.Println("Problem with coding msg ", err)
		return err
	}

	err = connection.writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: encodeMsg,
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
