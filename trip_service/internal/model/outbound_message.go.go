package model

type OutboundMessage struct {
	ID              string      `json:"id"`
	Source          string      `json:"source"`
	Type            string      `json:"type"`
	DataContentType string      `json:"datacontenttype"`
	Time            string      `json:"time"`
	Data            interface{} `json:"data"`
}
