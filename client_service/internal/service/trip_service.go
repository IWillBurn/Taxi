package service

import "offering_service/internal/socketlistener"

type TripService interface {
	CreateTrip(string, string) error
	CancelTrip(string, string) error
	GetTripStatus(string, string, *socketlistener.Publisher) error
}
