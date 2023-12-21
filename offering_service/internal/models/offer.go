package models

type Offer struct {
	From     LatLngLiteral `json:"from"`
	To       LatLngLiteral `json:"to"`
	ClientId string        `json:"client_id"`
	Price    Price         `json:"price"`
}
