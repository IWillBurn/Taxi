package publishers

import (
	"client_service/internal/models"
	"fmt"
	"sync"
)

type Publisher struct {
	Subscribers      sync.Map
	SubscribersMutex sync.Mutex
	Broadcast        chan models.SocketMessage
	Key              string
}

func (p *Publisher) Publish(to string, message interface{}) {
	fmt.Println("Published")
	fmt.Println(to)
	toMap := make(map[int]string)
	toMap[0] = to
	fmt.Println(toMap)
	p.Broadcast <- models.SocketMessage{
		To:      toMap,
		Message: message,
	}
}

func NewPublisher(key string) *Publisher {
	publisher := Publisher{
		Subscribers: sync.Map{},
		Broadcast:   make(chan models.SocketMessage),
		Key:         key,
	}
	return &publisher
}
