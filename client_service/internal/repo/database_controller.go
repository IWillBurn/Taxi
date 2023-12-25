package repo

type DataBaseController interface {
	GetTrips(userId string) ([]Trip, error)
	GetTripByTripId(tripId string) (Trip, error)
	GetTripByOfferId(offerId string) (Trip, error)
	AddTrip(trip Trip) error
	ChangeTripByTripId(tripId string, key string, value interface{}) error
	ChangeTripByOfferId(tripId string, key string, value interface{}) error
}
