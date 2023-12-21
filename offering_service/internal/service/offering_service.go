package service

import (
	"offering_service/internal/httpadapter/requests"
	"offering_service/internal/models"
)

type OfferingService interface {
	GetPrice(requests.CreateOfferRequest) models.Offer
}
