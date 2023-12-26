package main

import (
	socketrequests "client_service/internal/socketlistener/requests"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := "ws://localhost:53242/"
	log.Printf("Connecting to %s", u)
	header := http.Header{}
	header.Add("user_id", "client_1")
	c, _, err := websocket.DefaultDialer.Dial(u, header)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	// Горутина для чтения сообщений от сервера
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			fmt.Printf("Received message: %s\n", message)
		}
	}()
	tripId := "1"
	// Горутина для отправки сообщений на сервер
	go func() {
		defer close(done)
		for {
			select {
			case <-time.Tick(5 * time.Second):
				if tripId != "" {
					data := make(map[string]string)
					data["trip_id"] = tripId
					request := &socketrequests.SocketRequest{
						Key:  "status",
						Data: data,
					}
					err := c.WriteJSON(request)
					if err != nil {
						return
					}
				}
			}
		}
	}()

	select {
	case <-done:
	case <-interrupt:
		log.Println("interrupt")

		// Отправка сигнала завершения горутин
		err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("write close:", err)
			return
		}

		select {
		case <-done:
		case <-time.After(time.Second):
		}
	}
}
