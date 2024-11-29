package input

import "time"

type CreateStockInput struct {
	ProductID    int     `json:"product_id"`
	Quantity     int     `json:"quantity"`
	BasePrice    float64 `json:"base_price"`
	SellingPrice float64 `json:"selling_price"`
	PurchasePrice float64 `json:"purchase_price"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
}
