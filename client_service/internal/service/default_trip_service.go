package service

import (
	"fmt"
)

type DefaultTripService struct {
	Key string
}

func (tripService DefaultTripService) CreateTrip(string, string) error {
	return fmt.Errorf("error")
}

func (tripService DefaultTripService) CancelTrip(string, string) error {
	return fmt.Errorf("error")
}

func (tripService DefaultTripService) GetTripStatus(string, string) error {
	return fmt.Errorf("error")
}
