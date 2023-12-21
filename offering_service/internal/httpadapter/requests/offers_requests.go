package requests

import "offering_service/internal/models"

type CreateOfferRequest struct {
	From     models.LatLngLiteral `json:"from"`
	To       models.LatLngLiteral `json:"to"`
	ClientId string               `json:"client_id"`
}

type ParseOfferRequest struct {
	OfferId string `json:"offer_id"`
}
