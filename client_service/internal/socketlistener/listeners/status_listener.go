package listeners

import (
	"client_service/internal/socketlistener/publishers"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

type StatusListener struct {
	Clients   *map[string]*websocket.Conn
	ClientsMu *sync.Mutex
	Publisher *publishers.Publisher
}

func (l *StatusListener) HandleMessages() {
	for {
		msg := <-(*l.Publisher).Broadcast
		fmt.Println("Readed")
		l.ClientsMu.Lock()
		fmt.Println("Locked")
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
				fmt.Println("Sended")
				continue
			}

		}
		l.ClientsMu.Unlock()
		fmt.Println("Unlocked")
	}
}
