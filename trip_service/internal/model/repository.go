package model

import (
	"context"
)

type TripRepository interface {
	GetTrip(ctx context.Context, id int32) error
	//UpdateTrip(ctx context.Context, id int32) error
	//DeleteTrip(ctx context.Context, id int32) error
	//GetClientTrips(ctx context.Context, clientId int32) (*models.Trips, error)
}
