package handler

import (
	"api-kasirapp/formatter"
	"api-kasirapp/helper"
	"api-kasirapp/input"
	"api-kasirapp/service"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type customerHandler struct {
	customerService service.CustomerService
}

func NewCustomerHandler(customerService service.CustomerService) *customerHandler {
	return &customerHandler{customerService}
}

func (h *customerHandler) CreateCustomer(c *gin.Context) {
	var input input.CustomerInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Create customer failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newCustomer, err := h.customerService.CreateCustomer(input)
	if err != nil {
		response := helper.APIResponse("Create customer failed", http.StatusBadRequest, "error", gin.H{"message": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success create customer", http.StatusCreated, "success", formatter.FormatCustomer(newCustomer))
	c.JSON(http.StatusCreated, response)
}

func (h *customerHandler) GetCustomers(c *gin.Context) {
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

	customers, err := h.customerService.GetCustomers(limit, offset)
	if err != nil {
		response := helper.APIResponse("Get customers failed", http.StatusBadRequest, "error", gin.H{"message": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}

	totalCount, err := h.customerService.CountCustomers()
	if err != nil {
		response := helper.APIResponse("Get customers failed", http.StatusBadRequest, "error", gin.H{"message": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	formattedCustomers := formatter.FormatCustomers(customers)

	paginationMeta := gin.H{
		"total_data":   totalCount,
		"total_pages":  totalPages,
		"current_page": offset/limit + 1,
		"per_page":     limit,
	}

	response := helper.APIResponse("Success get customers", http.StatusOK, "success", gin.H{
		"data":       formattedCustomers,
		"pagination": paginationMeta,
	})
	c.JSON(http.StatusOK, response)
}

func (h *customerHandler) GetCustomerById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := helper.APIResponse("Invalid ID format", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	getCustomer, err := h.customerService.GetCustomerByID(id)
	if err != nil {
		if err.Error() == "record not found" {
			response := helper.APIResponse("Get customer failed", http.StatusNotFound, "error", nil)
			c.JSON(http.StatusNotFound, response)
			return
		}
		response := helper.APIResponse("Get customer failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success get customer", http.StatusOK, "success", formatter.FormatCustomer(getCustomer))
	c.JSON(http.StatusOK, response)
}

func (h *customerHandler) UpdateCustomer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := helper.APIResponse("Invalid ID format", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var input input.CustomerInput
	err = c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Update customer failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	updateCustomer, err := h.customerService.UpdateCustomer(id, input)
	if err != nil {
		response := helper.APIResponse("Update customer failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success update customer", http.StatusOK, "success", formatter.FormatCustomer(updateCustomer))
	c.JSON(http.StatusOK, response)
}

func (h *customerHandler) DeleteCustomer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := helper.APIResponse("Invalid ID format", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	deleteCustomer, err := h.customerService.DeleteCustomer(id)
	if err != nil {
		response := helper.APIResponse("Delete customer failed", http.StatusNotFound, "error", gin.H{"message": err.Error()})
		c.JSON(http.StatusNotFound, response)
		return
	}

	response := helper.APIResponse("Success delete customer", http.StatusOK, "success", formatter.FormatCustomer(deleteCustomer))
	c.JSON(http.StatusOK, response)
}

func (h *customerHandler) ExportCustomers(c *gin.Context) {
	// Panggil fungsi service untuk mengekspor data pelanggan ke file Excel
	file, err := h.customerService.ExportCustomersToXLS()
	if err != nil {
		// Jika terjadi error, kirimkan respons HTTP dengan status 500
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set header untuk file download
	c.Header("Content-Disposition", `attachment; filename="customers.xlsx"`)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	// Tulis file Excel langsung ke response writer
	if err := file.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write Excel file"})
	}
}

func (h *customerHandler) ImportCustomers(c *gin.Context) {
	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		response := helper.APIResponse("Failed to process file", http.StatusBadRequest, "error", gin.H{"message": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Save the file temporarily
	tempFilePath := "./tmp/" + file.Filename
	if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
		response := helper.APIResponse("Failed to save file", http.StatusInternalServerError, "error", gin.H{"message": err.Error()})
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	defer os.Remove(tempFilePath) // Clean up the temporary file

	// Import customers from the file
	customers, err := h.customerService.ImportCustomersFromXLS(tempFilePath)
	if err != nil {
		response := helper.APIResponse("Failed to import customers", http.StatusInternalServerError, "error", gin.H{"message": err.Error()})
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIResponse("Successfully imported customers", http.StatusOK, "success", customers)
	c.JSON(http.StatusOK, response)
}

