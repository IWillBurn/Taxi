package service

import (
	"offering_service/internal/models"
)

type DataBaseService interface {
	GetTrips(string) []models.OfferRecord
	GetTripById(string) (models.OfferRecord, error)
}
