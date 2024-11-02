package models

import "time"

type Discount struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Percentage float64   `json:"percentage"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
