package models

import "time"

type Customer struct {
	ID        int
	Name      string
	Address   string
	Phone     string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
