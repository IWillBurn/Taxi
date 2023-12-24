package models

type Offer struct {
	From     LatLngLiteral `json:"from"`
	To       LatLngLiteral `json:"to"`
	ClientId string        `json:"client_id"`
	Price    Price         `json:"price"`
}

type OfferRecord struct {
	Id       string        `json:"id"`
	OfferId  string        `json:"offer_id"`
	From     LatLngLiteral `json:"from"`
	To       LatLngLiteral `json:"to"`
	ClientId string        `json:"client_id"`
	Price    Price         `json:"price"`
	Status   string        `json:"status"`
}
