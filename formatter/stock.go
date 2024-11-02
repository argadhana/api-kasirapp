package formatter

import "api-kasirapp/models"

type StockFormatter struct {
	ID           int     `json:"id"`
	ProductName  string  `json:"product_name"`
	CodeProduct  string  `json:"code_product"`
	CategoryName string  `json:"category"`
	Stock        int     `json:"stock"`
	SellingPrice float64 `json:"selling_price"`
	BasePrice    float64 `json:"base_price"`
	Date         string  `json:"date"`
	BuyingPrice  float64 `json:"buying_price"`
	Amount       int     `json:"amount"`
	Information  string  `json:"information"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

func FormatStock(stock models.Stock) StockFormatter {
	formatter := StockFormatter{
		ID:           stock.ID,
		ProductName:  stock.ProductName,
		CodeProduct:  stock.CodeProduct,
		CategoryName: stock.CategoryName,
		Stock:        stock.Stock,
		SellingPrice: stock.SellingPrice,
		BasePrice:    stock.BasePrice,
		Date:         stock.Date,
		BuyingPrice:  stock.BuyingPrice,
		Amount:       stock.Amount,
		Information:  stock.Information,
		CreatedAt:    stock.CreatedAt.String(),
		UpdatedAt:    stock.UpdatedAt.String(),
	}
	return formatter
}
