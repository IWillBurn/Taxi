package models

type SocketMessage struct {
	From    string
	To      map[int]string
	Message interface{}
}
