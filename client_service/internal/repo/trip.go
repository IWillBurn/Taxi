package repo

import "client_service/internal/models"

type Trip struct {
	InitId   string               `json:"init_id"`
	TripId   string               `bson:"trip_id"`
	OfferId  string               `bson:"offer_id"`
	From     models.LatLngLiteral `bson:"from"`
	To       models.LatLngLiteral `bson:"to"`
	ClientId string               `bson:"client_id"`
	Price    models.Price         `bson:"price"`
	Status   string               `bson:"status"`
}
