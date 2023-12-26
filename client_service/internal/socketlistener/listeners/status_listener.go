package listeners

import (
	"client_service/internal/socketlistener/publishers"
	"github.com/gorilla/websocket"
	"sync"
)

type StatusListener struct {
	Clients   *sync.Map
	Publisher *publishers.Publisher
}

func (l *StatusListener) HandleMessages() {
	for {
		msg := <-(*l.Publisher).Broadcast
		for _, clientId := range msg.To {
			if clientId != "" {
				client, ok := l.Clients.Load(clientId)
				if !ok {
					continue
				}
				s, ok := client.(*websocket.Conn)
				err := s.WriteJSON(msg.Message)
				if err != nil {
					s.Close()
					l.Clients.Delete(clientId)
				}
				continue
			}

		}
	}
}
