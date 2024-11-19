package handler

import (
	"api-kasirapp/formatter"
	"api-kasirapp/helper"
	"api-kasirapp/input"
	"api-kasirapp/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type transactionHandler struct {
	transactionService service.OrderServices
}

func NewTransactionHandler(transactionService service.OrderServices) *transactionHandler {
	return &transactionHandler{transactionService}
}

func (h *transactionHandler) CreateTransaction(c *gin.Context) {
	// input from user
	var input input.TransactionInput

	// binding input
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Create transaction failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// create transaction
	newTransaction, cashReturn, err := h.transactionService.CreateTransactionWithCash(input)
	if err != nil {
		response := helper.APIResponse("Create transaction failed", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse(
		"Success create transaction",
		http.StatusCreated,
		"success",
		formatter.FormatTransaction(newTransaction, cashReturn),
	)
	c.JSON(http.StatusCreated, response)
}
