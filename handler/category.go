package handler

import (
	"api-kasirapp/formatter"
	"api-kasirapp/helper"
	"api-kasirapp/input"
	"api-kasirapp/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type categoryHandler struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) *categoryHandler {
	return &categoryHandler{
		categoryService,
	}
}

func (h *categoryHandler) CreateCategory(c *gin.Context) {
	var input input.CategoryInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Create category failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newCategory, err := h.categoryService.SaveCategory(input)
	if err != nil {
		response := helper.APIResponse("Create category failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success create category", http.StatusOK, "success", formatter.FormatCategory(newCategory))
	c.JSON(http.StatusOK, response)
}

func (h *categoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.categoryService.FindCategories()
	if err != nil {
		response := helper.APIResponse("Get categories failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatedCategories := formatter.FormatCategories(categories)

	response := helper.APIResponse("Success get categories", http.StatusOK, "success", formatedCategories)
	c.JSON(http.StatusOK, response)
}

func (h *categoryHandler) GetCategoryById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := helper.APIResponse("Invalid ID format", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	getCategory, err := h.categoryService.FindCategoryByID(id)
	if err != nil {
		response := helper.APIResponse("ID Not Found", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatedCategory := formatter.FormatCategory(getCategory)

	response := helper.APIResponse("Success get category", http.StatusOK, "success", formatedCategory)
	c.JSON(http.StatusOK, response)
}

func (h *categoryHandler) UpdateCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := helper.APIResponse("Invalid ID format", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var input input.CategoryInput
	err = c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Update category failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	updateCategory, err := h.categoryService.UpdateCategory(id, input)
	if err != nil {
		response := helper.APIResponse("Update category failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success update category", http.StatusOK, "success", formatter.FormatCategory(updateCategory))
	c.JSON(http.StatusOK, response)
}

func (h *categoryHandler) DeleteCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := helper.APIResponse("Invalid ID format", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	deleteCategory, err := h.categoryService.DeleteCategory(id)
	if err != nil {
		response := helper.APIResponse("Delete category failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success delete category", http.StatusOK, "success", formatter.FormatCategory(deleteCategory))
	c.JSON(http.StatusOK, response)
}

func (h *categoryHandler) GetCategoryProducts(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := helper.APIResponse("Invalid ID format", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	products, err := h.categoryService.GetCategoryProducts(id)
	if err != nil {
		response := helper.APIResponse("Get category products failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	categories, err := h.categoryService.FindCategoryByID(id)
	if err != nil {
		response := helper.APIResponse("ID Not Found", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return

	}

	formatedProducts := formatter.FormatCategoryProducts(categories, products)

	response := helper.APIResponse("Success get category products", http.StatusOK, "success", formatedProducts)
	c.JSON(http.StatusOK, response)
}

func (h *categoryHandler) GetProductsByCategoryName(c *gin.Context) {
	// Ambil parameter kategori dari request
	categoryName := c.Param("category_name")
	if categoryName == "" {
		response := helper.APIResponse("Category name is required", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	products, err := h.categoryService.GetProductsWithCategoryName(categoryName)
	if err != nil {
		response := helper.APIResponse("Get category products failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Panggil service
	category, err := h.categoryService.GetCategoryByName(categoryName)
	if err != nil {
		response := helper.APIResponse("Failed to get products by category", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formattedCategory := formatter.FormatCategoryProducts(category, products)

	// Format dan kembalikan response
	response := helper.APIResponse("Success get products by category", http.StatusOK, "success", formattedCategory)
	c.JSON(http.StatusOK, response)
}
