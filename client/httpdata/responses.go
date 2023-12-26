package httpdata

import "client/models"

type OfferResponse struct {
	Id       string               `json:"id"`
	From     models.LatLngLiteral `json:"from"`
	To       models.LatLngLiteral `json:"to"`
	ClientId string               `json:"client_id"`
	Price    models.Price         `json:"price"`
}

type TripResponse struct {
	TripId string `json:"trip_id"`
}

type SocketResponse struct {
	Status string `json:"status"`
}
