package socketlistener

import (
	"client_service/internal/service"
	"client_service/internal/socketlistener/publishers"
	socketrequests "client_service/internal/socketlistener/requests"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type Connector struct {
	TripService service.TripService
	Clients     map[string]*websocket.Conn
	ClientsMu   sync.Mutex
	Publishers  *map[string]*publishers.Publisher
}

func (c *Connector) HandleConnections(w http.ResponseWriter, r *http.Request) {
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
	c.ClientsMu.Lock()
	c.Clients[userId] = s
	c.ClientsMu.Unlock()
	for {
		request := &socketrequests.SocketRequest{}
		err = s.ReadJSON(request)
		if err != nil {
			fmt.Println(err)
			c.ClientsMu.Lock()
			delete(c.Clients, userId)
			c.ClientsMu.Unlock()
			return
		}
		if request.Key == "status" {
			fmt.Println("Connected")
			fmt.Println((*c.Publishers)[request.Key])
			fmt.Println("OK")
			err := (c.TripService).GetTripStatus(userId, request.Data["trip_id"], (*c.Publishers)[request.Key])
			if err != nil {
				return
			}
		}
	}
}
