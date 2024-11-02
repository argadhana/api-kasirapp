package models

import "time"

type User struct {
	ID           int
	Email        string
	PasswordHash string
	Name         string
	Phone        string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
