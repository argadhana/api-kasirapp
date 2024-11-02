package models

import "time"

type Stock struct {
	ID           int       `json:"id"`
	Date         string    `json:"date"`
	BuyingPrice  float64   `json:"buying_price"`
	Amount       int       `json:"amount"`
	Information  string    `json:"information"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
