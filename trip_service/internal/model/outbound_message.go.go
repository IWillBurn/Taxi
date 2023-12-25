package model

type OutboundMessage struct {
	ID              string      `json:"id"`
	Source          string      `json:"source"`
	Type            string      `json:"type"`
	DataContentType string      `json:"datacontenttype"`
	Time            string      `json:"time"`
	Data            interface{} `json:"data"`
}

type EventAcceptTrip struct {
	TripID   string `json:"trip_id"`
	DriverID string `json:"driver_id"`
}

type EventCancelTrip struct {
	TripID string `json:"trip_id"`
}

type EventCreatTrip struct {
	TripID  string        `json:"trip_id"`
	OfferID string        `json:"offer_id"`
	Price   Price         `json:"price"`
	Status  string        `json:"status"`
	From    LatLngLiteral `json:"from"`
	To      LatLngLiteral `json:"to"`
}

type EventStartTrip struct {
	TripID string `json:"trip_id"`
}

type EventEndTrip struct {
	TripID string `json:"trip_id"`
}
