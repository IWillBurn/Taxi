package httpdata

import "client/models"

type OfferRequest struct {
	From     models.LatLngLiteral `json:"from"`
	To       models.LatLngLiteral `json:"to"`
	ClientId string               `json:"client_id"`
}

type TripRequest struct {
	OfferId string `json:"offer_id"`
}

type SocketRequest struct {
	Key  string            `json:"key"`
	Data map[string]string `json:"data"`
}
