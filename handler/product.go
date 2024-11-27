package handler

import (
	"api-kasirapp/formatter"
	"api-kasirapp/helper"
	"api-kasirapp/input"
	"api-kasirapp/service"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type productHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *productHandler {
	return &productHandler{productService}
}

func (h *productHandler) CreateProduct(c *gin.Context) {
	var input input.ProductInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Create product failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newProduct, err := h.productService.CreateProduct(input)
	if err != nil {
		if err.Error() == "product code already exists" {
			response := helper.APIResponse("product code already exists", http.StatusConflict, "error", gin.H{"message": err.Error()})
			c.JSON(http.StatusConflict, response)
			return
		}
		response := helper.APIResponse("Create product failed", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success create product", http.StatusCreated, "success", formatter.FormatProduct(newProduct))
	c.JSON(http.StatusCreated, response)
}

func (h *productHandler) GetProducts(c *gin.Context) {
	products, err := h.productService.FindAll()
	if err != nil {
		response := helper.APIResponse("Get products failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success get products", http.StatusOK, "success", formatter.FormatProducts(products))
	c.JSON(http.StatusOK, response)
}

func (h *productHandler) GetProductById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := helper.APIResponse("Invalid ID format", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	getProduct, err := h.productService.FindProductByID(id)
	if err != nil {
		response := helper.APIResponse("Get product failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success get product", http.StatusOK, "success", formatter.FormatProduct(getProduct))
	c.JSON(http.StatusOK, response)
}

func (h *productHandler) UpdateProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := helper.APIResponse("Invalid ID format", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var input input.ProductInput
	err = c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Update product failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	updateProduct, err := h.productService.UpdateProduct(id, input)
	if err != nil {
		response := helper.APIResponse("Update product failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success update product", http.StatusOK, "success", formatter.FormatProduct(updateProduct))
	c.JSON(http.StatusOK, response)
}

func (h *productHandler) DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := helper.APIResponse("Invalid ID format", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	deleteProduct, err := h.productService.DeleteProduct(id)
	if err != nil {
		response := helper.APIResponse("Delete product failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success delete product", http.StatusOK, "success", formatter.FormatProduct(deleteProduct))
	c.JSON(http.StatusOK, response)
}

func (h *productHandler) ExportProducts(c *gin.Context) {
	file, err := h.productService.ExportProductsToXLS()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set headers for file download
	c.Header("Content-Disposition", `attachment; filename="products.xlsx"`)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	// Stream file to response
	file.Write(c.Writer)
}

func (h *productHandler) ImportProducts(c *gin.Context) {
	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		response := helper.APIResponse("Import products failed", http.StatusBadRequest, "error", gin.H{"message": "file not found"})
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Save the file to a temporary location
	filePath := "./temp_" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		response := helper.APIResponse("Import products failed", http.StatusInternalServerError, "error", gin.H{"message": "failed to save file"})
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	defer os.Remove(filePath) // Remove the file after processing

	// Call the import service
	importedProducts, err := h.productService.ImportProductsFromXLS(filePath)
	if err != nil {
		response := helper.APIResponse("Import products failed", http.StatusInternalServerError, "error", gin.H{"message": err.Error()})
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Respond with the list of imported products
	response := helper.APIResponse("Success import products", http.StatusOK, "success", gin.H{"products": importedProducts})
	c.JSON(http.StatusOK, response)
}
func (h *productHandler) UploadProductImage(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := helper.APIResponse("Invalid ID format", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Get the uploaded file
	file, err := c.FormFile("image")
	if err != nil {
		response := helper.APIResponse("Upload image failed", http.StatusBadRequest, "error", gin.H{"message": "file not found"})
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Define the target directory
	imageDir := "/var/www/images-product"

	// Ensure the directory exists
	if _, err := os.Stat(imageDir); os.IsNotExist(err) {
		if err := os.MkdirAll(imageDir, os.ModePerm); err != nil {
			response := helper.APIResponse("Failed to create image directory", http.StatusInternalServerError, "error", gin.H{"message": err.Error()})
			c.JSON(http.StatusInternalServerError, response)
			return
		}
	}

	// Save the file to the directory
	filePath := fmt.Sprintf("%s/%s", imageDir, file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		response := helper.APIResponse("Upload image failed", http.StatusInternalServerError, "error", gin.H{"message": "failed to save file"})
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Generate the public image URL
	publicURL := fmt.Sprintf("/images/%s", file.Filename)

	// Call the service to update the product with the image URL
	updatedProduct, err := h.productService.SaveProductImage(id, publicURL)
	if err != nil {
		response := helper.APIResponse("Upload image failed", http.StatusInternalServerError, "error", gin.H{"message": err.Error()})
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Respond with the public image URL
	response := helper.APIResponse("Success upload image", http.StatusOK, "success", gin.H{
		"image_url": publicURL,
		"product":   formatter.FormatProduct(updatedProduct),
	})
	c.JSON(http.StatusOK, response)
}
