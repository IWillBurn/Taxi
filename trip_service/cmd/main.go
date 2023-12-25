package main

import (
	"context"
	"flag"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"trip_service/internal/app"
	"trip_service/internal/app/config"
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
	//create_trip, err := os.ReadFile("./create_trip.json")
	//go func() {
	//	index := 200
	//	writer := kafka.NewWriter(kafka.WriterConfig{
	//		Brokers: []string{"kafka:29092"},
	//		Topic:   "trip_inbound",
	//	})
	//
	//	defer writer.Close()
	//	for {
	//		err = writer.WriteMessages(context.Background(), kafka.Message{Key: []byte(strconv.Itoa(index)), Value: create_trip})
	//		if err != nil {
	//			log.Println("Problem with writing message ", err)
	//		}
	//		time.Sleep(3 * time.Second)
	//		index += 1
	//	}
	//}()

	if err := application.Run(context.Background()); err != nil {
		log.Fatal("Error on server", zap.Error(err))
	}

	<-done
}
