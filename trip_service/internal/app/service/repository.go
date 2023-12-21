package service

import (
	"context"
	"trip_service/pkg/models"
)

type TripRepository interface {
	GetTrip(ctx context.Context, id int32) (*models.Trip, error)
	UpdateTrip(ctx context.Context, id int32) error
	DeleteTrip(ctx context.Context, id int32) error
	GetClientTrips(ctx context.Context, clientId int32) (*models.Trips, error)
}
