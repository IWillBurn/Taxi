package models

type Price struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}
