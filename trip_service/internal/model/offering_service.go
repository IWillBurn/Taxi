package model

type LatLngLiteral struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
type Price struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type Offer struct {
	Id       string        `json:"id"`
	From     LatLngLiteral `json:"from"`
	To       LatLngLiteral `json:"to"`
	ClientId string        `json:"client_id"`
	Price    Price         `json:"price"`
}
