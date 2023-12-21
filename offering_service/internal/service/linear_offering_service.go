package service

import (
	"math"
	"offering_service/internal/httpadapter/requests"
	"offering_service/internal/models"
)

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

type LinearOfferingService struct {
	LinearCost   float64
	BaseCost     float64
	PlanetRadius float64 // такси на других планетах :D
}

func (offering LinearOfferingService) GetPrice(request requests.CreateOfferRequest) models.Offer {

	lat1 := degreesToRadians(request.From.Lat)
	lon1 := degreesToRadians(request.From.Lng)
	lat2 := degreesToRadians(request.To.Lat)
	lon2 := degreesToRadians(request.To.Lng)

	orderLen := offering.PlanetRadius * math.Acos(math.Sin(lat1)*math.Sin(lat2)+math.Cos(lat1)*math.Cos(lat2)*math.Cos(lon2-lon1)) * 1000

	cost := offering.LinearCost*orderLen + offering.BaseCost
	return models.Offer{
		From:     request.From,
		To:       request.To,
		ClientId: request.ClientId,
		Price:    models.Price{Amount: cost, Currency: "RUB"},
	}
}
