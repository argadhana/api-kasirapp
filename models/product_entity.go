package models

import (
	"time"
)

type Product struct {
	ID           int
	Name         string
	ProductType  string
	ProductFileName	 string
	BasePrice    float64
	SellingPrice float64
	Stock        int
	CodeProduct  string
	CategoryID   int
	MinimumStock int
	Shelf        string
	Weight       int
	Discount     int
	Information  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
