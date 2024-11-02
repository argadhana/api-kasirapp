package models

import "time"

type Category struct {
	ID        int
	Name      string
	Product  []Product
	CreatedAt time.Time
	UpdatedAt time.Time
}