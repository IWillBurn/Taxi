package models

type OutboundMessage struct {
	Id              string      `json:"id"`
	Source          string      `json:"source"`
	Type            string      `json:"type"`
	DataContentType string      `json:"datacontenttype"`
	Time            string      `json:"time"`
	Data            interface{} `json:"data"`
}

type EventTripStatusUpdate struct {
	TripId string `json:"trip_id"`
}

type EventCreateTrip struct {
	TripId  string        `json:"trip_id"`
	OfferId string        `json:"offer_id"`
	Price   Price         `json:"price"`
	Status  string        `json:"status"`
	From    LatLngLiteral `json:"from"`
	To      LatLngLiteral `json:"to"`
}
