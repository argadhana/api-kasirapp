package models

import "time"

type Supplier struct {
	ID        int
	Name      string
	Address   string
	Email     string
	Phone     string
	Code      int
	CreatedAt time.Time
	UpdatedAt time.Time
}
