package tests

import (
	"github.com/stretchr/testify/assert"
	"offering_service/internal/models"
	"offering_service/internal/service"
	"testing"
)

func TestJWTSigningService_Coding(t *testing.T) {
	key := "secret_key"
	jwtService := service.JWTSigningService{Key: key}
	cases := []struct {
		name              string
		input             models.Offer
		expectEncodeError error
		expectDecodeError error
	}{
		{
			name: "test_encoding_1",
			input: models.Offer{
				From:     models.LatLngLiteral{0, 0},
				To:       models.LatLngLiteral{0, 0},
				ClientId: "client_1",
				Price:    models.Price{100, "RUB"},
			},
			expectEncodeError: nil,
			expectDecodeError: nil,
		},
		{
			name: "test_encoding_2",
			input: models.Offer{
				From:     models.LatLngLiteral{0, 0},
				To:       models.LatLngLiteral{10, 10},
				ClientId: "client_2",
				Price:    models.Price{1000, "EUR"},
			},
			expectEncodeError: nil,
			expectDecodeError: nil,
		},
	}
	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			encode, err := jwtService.Encode(cs.input)
			assert.Equal(t, cs.expectEncodeError, err)
			decode, err := jwtService.Decode(encode)
			assert.Equal(t, cs.expectDecodeError, err)
			assert.Equal(t, cs.input, decode)
		})
	}
}
