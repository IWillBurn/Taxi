package model

type InboundMessage struct {
	Id              string                 `json:"id"`
	Source          string                 `json:"source"`
	Type            string                 `json:"type"`
	DataContentType string                 `json:"datacontenttype"`
	Time            string                 `json:"time"`
	Data            map[string]interface{} `json:"data"`
}

type AcceptTrip struct {
	TripID   string `json:"trip_id"`
	DriverID string `json:"driver_id"`
}

type CancelTrip struct {
	TripID string `json:"trip_id"`
	Reason string `json:"driver_id"`
}

type CreatTrip struct {
	OfferID string `json:"offer_id"`
}

type StartTrip struct {
	TripID string `json:"trip_id"`
}

type EndTrip struct {
	TripID string `json:"trip_id"`
}
