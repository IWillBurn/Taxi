package models

type InboundMessage struct {
	Id              string                 `json:"id"`
	Source          string                 `json:"source"`
	Type            string                 `json:"type"`
	DataContentType string                 `json:"datacontenttype"`
	Time            string                 `json:"time"`
	Data            map[string]interface{} `json:"data"`
}

type CreateData struct {
	OfferId string `json:"offer_id"`
}

type CancelData struct {
	TripId string `json:"trip_id"`
	Reason string `json:"reason"`
}
