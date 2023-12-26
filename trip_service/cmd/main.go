package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"trip_service/internal/app"
	"trip_service/internal/app/config"
	"trip_service/internal/model"
)

func getConfigPath() string {
	var configPath string

	flag.StringVar(&configPath, "c", "./.config/trip_service.yaml", "path to config file")
	flag.Parse()

	return configPath
}

func main() {
	cfg, err := config.NewConfig(getConfigPath())
	if err != nil {
		log.Fatalln("Error on read config", err)
	}
	application := app.NewApp(context.Background(), cfg)

	done := make(chan os.Signal, 1)

	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT)
	// Driver service test
	//TODO delete this code
	go func() {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{"kafka:29092"},
			Topic:   "trip_outbound",
		})
		defer reader.Close()
		writer := kafka.NewWriter(kafka.WriterConfig{
			Brokers: []string{"kafka:29092"},
			Topic:   "trip_inbound",
		})
		defer writer.Close()
		for {
			msg, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Println(err)
			}
			var responseMessage model.OutboundMessage
			err = json.NewDecoder(strings.NewReader(string(msg.Value))).Decode(&responseMessage)
			//if responseMessage.Type == "trip.event.created" {
			//	data, _ := json.Marshal(responseMessage.Data)
			//	var responseSub model.EventCreatTrip
			//	_ = json.Unmarshal(data, &responseSub)
			//	encodeMsg, _ := json.Marshal(model.OutboundMessage{
			//		ID:              responseMessage.ID,
			//		Source:          "/driver",
			//		DataContentType: responseMessage.DataContentType,
			//		Type:            "trip.command.accept",
			//		Time:            responseMessage.Time,
			//		Data: model.AcceptTrip{
			//			TripID:   responseSub.TripID,
			//			DriverID: "1",
			//		}})
			//	_ = writer.WriteMessages(context.Background(), kafka.Message{Value: encodeMsg})
			//	time.Sleep(10 * time.Second)
			//	encodeMsg, _ = json.Marshal(model.OutboundMessage{
			//		ID:              responseMessage.ID,
			//		Source:          "/driver",
			//		DataContentType: responseMessage.DataContentType,
			//		Type:            "trip.command.start",
			//		Time:            responseMessage.Time,
			//		Data: model.StartTrip{
			//			TripID: responseSub.TripID,
			//		}})
			//	_ = writer.WriteMessages(context.Background(), kafka.Message{Value: encodeMsg})
			//	time.Sleep(10 * time.Second)
			//	encodeMsg, _ = json.Marshal(model.OutboundMessage{
			//		ID:              responseMessage.ID,
			//		Source:          "/driver",
			//		DataContentType: responseMessage.DataContentType,
			//		Type:            "trip.command.end",
			//		Time:            responseMessage.Time,
			//		Data: model.EndTrip{
			//			TripID: responseSub.TripID,
			//		}})
			//	_ = writer.WriteMessages(context.Background(), kafka.Message{Value: encodeMsg})
			//	time.Sleep(10 * time.Second)
			//}

		}
	}()

	if err := application.Run(context.Background()); err != nil {
		log.Fatal("Error on server", zap.Error(err))
	}

	<-done
}
