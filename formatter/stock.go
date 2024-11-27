package formatter

import (
	"api-kasirapp/models"
	"fmt"
)

type StockResponse struct {
	ID           int              `json:"id"`
	ProductID    int              `json:"product_id"`
	Product      ProductFormatter `json:"product"`
	Quantity     int              `json:"quantity"`
	BasePrice    string           `json:"base_price"`
	SellingPrice string           `json:"selling_price"`
	Date         string           `json:"date"`
	Description  string           `json:"description"`
}

func FormatStockResponse(stock models.Stock) StockResponse {
	return StockResponse{
		ID:           stock.ID,
		ProductID:    stock.ProductID,
		Product:      FormatProduct(stock.Product),
		Quantity:     stock.Quantity,
		BasePrice:    fmt.Sprintf("%.2f", stock.BasePrice),
		SellingPrice: fmt.Sprintf("%.2f", stock.SellingPrice),
		Date:         stock.Date.Format("2006-01-02"),
		Description:  stock.Description,
	}
}

func FormatStocks(stocks []models.Stock) []StockResponse {
	var stocksResponse []StockResponse

	for _, stock := range stocks {
		formattedStock := FormatStockResponse(stock)
		stocksResponse = append(stocksResponse, formattedStock)
	}

	return stocksResponse
}
