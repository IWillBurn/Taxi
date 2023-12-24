package socketlistener

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"offering_service/internal/service"
	socketrequests "offering_service/internal/socketlistener/requests"
	"sync"
)

type Connector struct {
	TripService *service.TripService
	Clients     map[string]*websocket.Conn
	ClientsMu   sync.Mutex
	Publishers  *map[string]*Publisher
}

func (c *Connector) HandleConnections(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("user_id")

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()

	c.ClientsMu.Lock()
	c.Clients[userId] = ws
	c.ClientsMu.Unlock()

	for {
		request := &socketrequests.SocketRequest{}
		err := ws.ReadJSON(request)
		if err != nil {
			fmt.Println(err)
			c.ClientsMu.Lock()
			delete(c.Clients, userId)
			c.ClientsMu.Unlock()
			break
		}
		if request.Key == "status" {
			err := (*c.TripService).GetTripStatus(userId, request.Data["trip_id"], (*c.Publishers)["request.Key"])
			if err != nil {
				return
			}
		}
	}
}
