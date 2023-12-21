package test

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/segmentio/kafka-go"
)

var async = flag.Bool("a", false, "use async")

func main() {
	flag.Parse()

	ctx := context.Background()

	logger := log.Default()

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:     []string{"kafka:9092"},
		Topic:       "trip_service",
		Async:       *async,
		Logger:      kafka.LoggerFunc(logger.Printf),
		ErrorLogger: kafka.LoggerFunc(logger.Printf),
		BatchSize:   2000,
	})
	defer writer.Close()

	acceptTrip, _ := ioutil.ReadFile("/accept_trip.json")

	for i := 0; i < 524_288; i++ {
		err := writer.WriteMessages(ctx, kafka.Message{Key: []byte(strconv.Itoa(i)), Value: acceptTrip})
		if err != nil {
			log.Fatal(err)
		}
	}
}
