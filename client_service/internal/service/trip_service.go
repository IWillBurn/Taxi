package service

import (
	"client_service/internal/socketlistener/publishers"
)

type TripService interface {
	CreateTrip(offerId string) error
	CancelTrip(tripId string, reason string) error
	GetTripStatus(clientId string, TripId string, publisher *publishers.Publisher) error
}
