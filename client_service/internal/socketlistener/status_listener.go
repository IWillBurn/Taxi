package socketlistener

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

type StatusListener struct {
	Clients   *map[string]*websocket.Conn
	ClientsMu *sync.Mutex
	Publisher *Publisher
}

func (l *StatusListener) HandleMessages() {
	for {
		msg := <-(*l.Publisher).Broadcast
		l.ClientsMu.Lock()
		for _, clientId := range msg.To {
			if clientId != "" {
				client, ok := (*l.Clients)[clientId]
				if !ok {
					continue
				}
				err := client.WriteJSON(msg.Message)
				if err != nil {
					fmt.Println(err)
					client.Close()
					delete(*l.Clients, clientId)
				}
				continue
			}

		}
		l.ClientsMu.Unlock()
	}
}
