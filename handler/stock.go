// handler/stock_handler.go
package handler

import (
	"api-kasirapp/input"
	"api-kasirapp/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StockHandler struct {
	stockService service.StockService
}

func NewStockHandler(stockService service.StockService) *StockHandler {
	return &StockHandler{stockService}
}

func (h *StockHandler) AddStock(c *gin.Context) {
	var input input.CreateStockInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newStock, err := h.stockService.AddStock(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": newStock})
}

func (h *StockHandler) GetStocks(c *gin.Context) {
	stocks, err := h.stockService.GetStocks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stocks})
}

func (h *StockHandler) GetStocksByProductID(c *gin.Context) {
	productID := c.Param("productID")

	productIDInt, err := strconv.Atoi(productID)

	stocks, err := h.stockService.GetStocksByProductID(productIDInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stocks})
}
