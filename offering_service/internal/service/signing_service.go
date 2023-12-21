package service

import "offering_service/internal/models"

type SigningService interface {
	Encode(data interface{}) (string, error)
	Decode(tokenString string) (models.Offer, error)
}
