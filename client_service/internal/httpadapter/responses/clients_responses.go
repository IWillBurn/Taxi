package responses

import "client_service/internal/models"

type OfferResponse struct {
	Id       string               `json:"id"`
	From     models.LatLngLiteral `json:"from"`
	To       models.LatLngLiteral `json:"to"`
	ClientId string               `json:"client_id"`
	Price    models.Price         `json:"price"`
}
