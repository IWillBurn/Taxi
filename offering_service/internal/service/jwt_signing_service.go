package service

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"offering_service/internal/models"
	"time"
)

type JWTSigningService struct {
	Key string
}

func (jwtService JWTSigningService) Encode(data interface{}) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["data"] = data
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // Токен действителен в течение 1 часа

	secretKey := []byte(jwtService.Key)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("signing error")
	}

	return tokenString, nil
}

func (jwtService JWTSigningService) Decode(tokenString string) (models.Offer, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("incorrect signature method: %v", token.Header["alg"])
		}
		return []byte(jwtService.Key), nil
	})

	if err != nil {
		return models.Offer{}, fmt.Errorf("parsing_error")
	}

	if token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return models.Offer{}, fmt.Errorf("claims_error")
		}

		data := claims["data"].(map[string]interface{})

		jsonData, err := json.Marshal(data)
		if err != nil {
			return models.Offer{}, fmt.Errorf("marshal_error")
		}

		var offer models.Offer
		err = json.Unmarshal(jsonData, &offer)
		if err != nil {
			return models.Offer{}, fmt.Errorf("unmarshal_error")
		}
		return offer, nil

	} else {
		return models.Offer{}, fmt.Errorf("invalid_token")
	}

}
