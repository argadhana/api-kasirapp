package handler

import (
	"api-kasirapp/formatter"
	"api-kasirapp/helper"
	"api-kasirapp/input"
	"api-kasirapp/service"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type supplierHandler struct {
	supplierService service.SupplierService
}

func NewSupplierHandler(supplierService service.SupplierService) *supplierHandler {
	return &supplierHandler{supplierService}
}

func (h *supplierHandler) CreateSupplier(c *gin.Context) {
	var input input.SupplierInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Create supplier failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newSupplier, err := h.supplierService.CreateSupplier(input)
	if err != nil {
		response := helper.APIResponse("Create supplier failed", http.StatusBadRequest, "error", gin.H{"message": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success create supplier", http.StatusCreated, "success", formatter.FormatSupplier(newSupplier))
	c.JSON(http.StatusCreated, response)
}

func (h *supplierHandler) GetSupplierById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := helper.APIResponse("Invalid ID format", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	getSupplier, err := h.supplierService.GetByID(id)
	if err != nil {
		response := helper.APIResponse("Get supplier failed", http.StatusNotFound, "error", gin.H{"message": err.Error()})
		c.JSON(http.StatusNotFound, response)
		return
	}

	response := helper.APIResponse("Success get supplier", http.StatusOK, "success", formatter.FormatSupplier(getSupplier))
	c.JSON(http.StatusOK, response)
}

func (h *supplierHandler) GetSuppliers(c *gin.Context) {
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 4
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	suppliers, err := h.supplierService.GetAll(limit, offset)
	if err != nil {
		response := helper.APIResponse("Get suppliers failed", http.StatusBadRequest, "error", gin.H{"message": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success get suppliers", http.StatusOK, "success", formatter.FormatSuppliers(suppliers))
	c.JSON(http.StatusOK, response)
}

func (h *supplierHandler) UpdateSupplier(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := helper.APIResponse("Invalid ID format", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var input input.SupplierInput

	err = c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Update supplier failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	updateSupplier, err := h.supplierService.Update(id, input)
	if err != nil {
		response := helper.APIResponse("Update supplier failed", http.StatusBadRequest, "error", gin.H{"message": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success update supplier", http.StatusOK, "success", formatter.FormatSupplier(updateSupplier))
	c.JSON(http.StatusOK, response)
}

func (h *supplierHandler) DeleteSupplier(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := helper.APIResponse("Invalid ID format", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	deleteSupplier, err := h.supplierService.Delete(id)
	if err != nil {
		response := helper.APIResponse("Delete supplier failed", http.StatusBadRequest, "error", gin.H{"message": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success delete supplier", http.StatusOK, "success", formatter.FormatSupplier(deleteSupplier))
	c.JSON(http.StatusOK, response)
}

func (h *supplierHandler) ExportSuppliers(c *gin.Context) {
	file, err := h.supplierService.ExportSuppliersToXLS()
	if err != nil {
		response := helper.APIResponse("Export suppliers failed", http.StatusInternalServerError, "error", gin.H{"message": err.Error()})
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", `attachment; filename="suppliers.xlsx"`)
	c.Header("File-Name", "suppliers.xlsx")
	file.Write(c.Writer)
}

func (h *supplierHandler) ImportSuppliers(c *gin.Context) {
	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		response := helper.APIResponse("Failed to get file", http.StatusBadRequest, "error", gin.H{"message": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Save the file to a temporary location
	filePath := "./tmp/%s" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		response := helper.APIResponse("Failed to save file", http.StatusInternalServerError, "error", gin.H{"message": err.Error()})
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	defer os.Remove(filePath) // Cleanup temporary file

	// Import suppliers from the file
	importedSuppliers, err := h.supplierService.ImportSuppliersFromXLS(filePath)
	if err != nil {
		response := helper.APIResponse("Failed to import suppliers", http.StatusInternalServerError, "error", gin.H{"message": err.Error()})
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIResponse("Successfully imported suppliers", http.StatusOK, "success", importedSuppliers)
	c.JSON(http.StatusOK, response)
}
