package service

import (
	"offering_service/internal/models"
)

type MongoService struct {
	Key string
}

func (mongoService MongoService) GetTrips(string) []models.OfferRecord {
	return []models.OfferRecord{}
}
