package repository

import (
	"context"
)

type Repository struct {
}

func (repo Repository) GetTrip(ctx context.Context, id int32) error {
	return nil
}

func NewRepository() *Repository {
	return &Repository{}
}

//func (repo *Repository) UpdateTrip(ctx context.Context, id int32) error{}
//func (repo *Repository) DeleteTrip(ctx context.Context, id int32) error{}
//func (repo *Repository) GetClientTrips(ctx context.Context, clientId int32) (*models.Trips, error){}
