package socketlistener

import (
	"client_service/internal/service"
	"client_service/internal/socketlistener/publishers"
	socketrequests "client_service/internal/socketlistener/requests"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type Connector struct {
	TripService service.TripService
	Clients     *sync.Map
	Publishers  *map[string]*publishers.Publisher
}

func (c *Connector) HandleConnections(w http.ResponseWriter, r *http.Request) {
	log.Printf("CONNECTION")
	userId := r.Header.Get("user_id")
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	s, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Clients.Store(userId, s)
	log.Println("ADD")
	log.Println(c.Clients)
	for {
		request := &socketrequests.SocketRequest{}
		err = s.ReadJSON(request)
		if request.Key == "status" {
			fmt.Println("Connected")
			fmt.Println((*c.Publishers)[request.Key])
			fmt.Println("OK")
			fmt.Println("userId")
			fmt.Println(userId)
			fmt.Println("trip_id")
			fmt.Println(request.Data["trip_id"])
			err := (c.TripService).GetTripStatus(userId, request.Data["trip_id"], (*c.Publishers)[request.Key])
			if err != nil {
				return
			}
		}
	}
}
