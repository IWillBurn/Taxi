package main

import (
	"bytes"
	"client/httpdata"
	"client/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func GetOfferId(userId string, client *http.Client) (string, error) {
	request := httpdata.OfferRequest{
		From: models.LatLngLiteral{
			Lat: 0,
			Lng: 0,
		},
		To: models.LatLngLiteral{
			Lat: 0,
			Lng: 0,
		},
		ClientId: userId,
	}

	requestJson, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Ошибка маршалинга данных в JSON:", err)
		return "", err
	}

	offerURL := "http://localhost:9090/offers"
	resp, err := client.Post(offerURL, "application/json", bytes.NewBuffer(requestJson))
	if err != nil {
		fmt.Println("Ошибка пост:", err)
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка чтения:", err)
		return "", err
	}
	var response httpdata.OfferResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Ошибка анмаршала:", err)
		return "", err
	}
	return response.Id, nil
}

func GetTripId(offerId string, userId string, client *http.Client) (string, error) {
	request := httpdata.TripRequest{
		OfferId: offerId,
	}

	requestJson, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Ошибка маршалинга данных в JSON:", err)
		return "", err
	}

	url := "http://localhost:8080/trips"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestJson))
	if err != nil {
		fmt.Println("Ошибка создания запроса:", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("user_id", userId)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Ошибка:", err)
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return "", err
	}
	var response httpdata.TripResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return "", err
	}
	return response.TripId, nil
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	client := &http.Client{}
	userId := "client_1"
	offerId, err := GetOfferId(userId, client)
	if err != nil {
		fmt.Println("Ошибка:", err)
	}
	tripId, err := GetTripId(offerId, userId, client)
	if err != nil {
		fmt.Println("Ошибка:", err)
	}

	// Создание сокета
	u := "ws://localhost:53242/"
	log.Printf("Connecting to %s", u)
	header := http.Header{}
	header.Add("user_id", userId)
	c, _, err := websocket.DefaultDialer.Dial(u, header)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	// read
	go func() {
		defer close(done)
		for {
			message := &httpdata.SocketResponse{}
			err := c.ReadJSON(message)
			if err != nil {
				log.Println("read:", err)
				return
			}
			fmt.Println(message.Status)
			if message.Status == "ENDED" {
				break
			}
		}
	}()

	// write
	go func() {
		defer close(done)
		for {
			select {
			case <-time.Tick(5 * time.Second):
				if tripId != "" {
					data := make(map[string]string)
					data["trip_id"] = tripId
					request := &httpdata.SocketRequest{
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
