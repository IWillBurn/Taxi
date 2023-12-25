package model

type Trip struct {
	Id           string
	OfferId      string
	DriverId     string
	CurrentStage string
}
type Trips []*Trip
