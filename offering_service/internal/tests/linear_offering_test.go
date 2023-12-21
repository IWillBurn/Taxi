package tests

import (
	"github.com/stretchr/testify/assert"
	"offering_service/internal/httpadapter/requests"
	"offering_service/internal/models"
	"offering_service/internal/service"
	"testing"
)

func TestLinearOfferingService_GetPrice(t *testing.T) {
	linearOfferingService := service.LinearOfferingService{LinearCost: 1, BaseCost: 100, PlanetRadius: 6370}
	cases := []struct {
		name   string
		input  requests.CreateOfferRequest
		expect models.Offer
	}{
		{
			name: "test_get_price_1",
			input: requests.CreateOfferRequest{
				From:     models.LatLngLiteral{0, 0},
				To:       models.LatLngLiteral{0, 0},
				ClientId: "client_1",
			},
			expect: models.Offer{
				From:     models.LatLngLiteral{0, 0},
				To:       models.LatLngLiteral{0, 0},
				ClientId: "client_1",
				Price:    models.Price{100, "RUB"},
			},
		},
		{
			name: "test_get_price_2",
			input: requests.CreateOfferRequest{
				From:     models.LatLngLiteral{0, 0},
				To:       models.LatLngLiteral{10, 10},
				ClientId: "client_1",
			},
			expect: models.Offer{
				From:     models.LatLngLiteral{0, 0},
				To:       models.LatLngLiteral{10, 10},
				ClientId: "client_1",
				Price:    models.Price{1568374.3598817938, "RUB"},
			},
		},
	}
	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			res := linearOfferingService.GetPrice(cs.input)
			assert.Equal(t, cs.expect, res)
		})
	}
}
