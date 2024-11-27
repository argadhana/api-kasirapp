// handler/stock_handler.go
package handler

import (
	"api-kasirapp/helper"
	"api-kasirapp/input"
	"api-kasirapp/service"
	"math"
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
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	stocks, err := h.stockService.GetStocks(limit, offset)
	if err != nil {
		response := helper.APIResponse("Get stocks failed", http.StatusBadRequest, "error", gin.H{"message": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}

	totalCount, err := h.stockService.CountStocks()
	if err != nil {
		response := helper.APIResponse("Get stocks failed", http.StatusBadRequest, "error", gin.H{"message": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	paginationMeta := gin.H{
		"total_data":   totalCount,
		"total_pages":  totalPages,
		"current_page": offset/limit + 1,
		"per_page":     limit,
	}

	response := helper.APIResponse("Success get stocks", http.StatusOK, "success", gin.H{
		"data":       stocks, // Format if needed using a formatter
		"pagination": paginationMeta,
	})
	c.JSON(http.StatusOK, response)
}

func (h *StockHandler) GetStocksByProductID(c *gin.Context) {
	productID := c.Param("productID")

	productIDInt, err := strconv.Atoi(productID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stocks, err := h.stockService.GetStocksByProductID(productIDInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stocks})
}
