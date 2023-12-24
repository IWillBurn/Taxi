package socketlistener

import (
	"sync"
)

type Message struct {
	From    string
	To      []string
	Message string
}

type Publisher struct {
	Subscribers      map[string]bool
	SubscribersMutex sync.Mutex
	Broadcast        chan Message
	Key              string
}

func NewPublisher(key string) *Publisher {
	publisher := Publisher{
		Subscribers: make(map[string]bool),
		Broadcast:   make(chan Message),
		Key:         key,
	}
	return &publisher
}
