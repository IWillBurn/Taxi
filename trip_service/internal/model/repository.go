package model

import (
	"context"
)

type ParamsStarted struct {
	Id           string
	OfferId      string
	DriverId     string
	CurrentStage string
	Offset       uint64
	Limit        uint64
}
type ParamsFinished struct {
	ParamsStarted
}
type TripStarted struct {
	Id           string
	OfferId      string
	DriverId     string
	CurrentStage string
}
type TripFinished struct {
	TripStarted
	Successful bool
	Reason     string
}

type TripRepository interface {
	GetStarted(ctx context.Context, params *ParamsStarted) ([]TripStarted, error)
	CreateStarted(ctx context.Context, trip *TripStarted) (string, error)
	UpdateStarted(ctx context.Context, params *ParamsStarted, value *TripStarted) error
	DeleteStarted(ctx context.Context, params *ParamsStarted) error

	GetFinished(ctx context.Context, params *ParamsFinished) ([]TripFinished, error)
	CreateFinished(ctx context.Context, trip *TripFinished) (string, error)
	UpdateFinished(ctx context.Context, params *ParamsFinished, value *TripFinished) error
}
