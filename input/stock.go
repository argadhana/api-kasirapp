package input

import "time"

type CreateStockInput struct {
	ProductID    int       `json:"product_id" binding:"required"`
	Quantity     int       `json:"quantity" binding:"required"`
	BasePrice    float64   `json:"base_price" binding:"required"`
	SellingPrice float64   `json:"selling_price" binding:"required"`
	Date         time.Time `json:"date" binding:"required"`
	Description  string    `json:"description"`
}
