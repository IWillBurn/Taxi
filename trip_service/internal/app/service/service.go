package service

import (
	"log"
	"trip_service/internal/model"
)

type TripService struct {
	repo *TripRepository
}

func NewService(
	repo *TripRepository,
) *TripService {
	return &TripService{
		repo: repo,
	}
}

func (s *TripService) AcceptTrip(msg model.AcceptTrip) {
	log.Println("TEST Accept trip", msg)
}
func (s *TripService) CreatTrip(msg model.CreatTrip)   {}
func (s *TripService) CancelTrip(msg model.CancelTrip) {}
func (s *TripService) EndTrip(msg model.EndTrip)       {}
func (s *TripService) StartTrip(msg model.StartTrip)   {}

//func (s *TripService) getTrip(ctx context.Context, id int32) (*models.Trip, error) {
//
//	return s.repo.GetTrip(ctx, id)
//}
//
//func (s *TripService) updateTrip(ctx context.Context, id int32) error {
//	return s.repo.UpdateTrip(ctx, id)
//}
//func (s *TripService) deleteTrip(ctx context.Context, id int32) error {
//	return s.repo.DeleteTrip(ctx, id)
//}
//
//func (s *TripService) getClientTrips(ctx context.Context, clientId int32) (*models.Trips, error) {
//	return s.repo.GetClientTrips(ctx, clientId)
//}
