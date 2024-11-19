package formatter

import "api-kasirapp/models"

type CategoryFormatter struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func FormatCategory(category models.Category) CategoryFormatter {
	formatter := CategoryFormatter{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: category.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return formatter
}

type ProductWithTheCategoryFormatter struct {
	Name string `json:"name"`
}

type CategoryProductFormatter struct {
	CategoryID   int                `json:"category_id"`
	CategoryName string             `json:"category_name"`
	Products     []ProductFormatter `json:"products"`
}

func FormatCategories(categories []models.Category) []CategoryFormatter {
	var categoriesFormatter []CategoryFormatter

	for _, category := range categories {
		formatter := FormatCategory(category)
		categoriesFormatter = append(categoriesFormatter, formatter)
	}

	return categoriesFormatter
}

func FormatCategoryProducts(category models.Category, products []models.Product) CategoryProductFormatter {
	var productList []ProductFormatter
	for _, product := range products {
		formattedProduct := FormatProduct(product)
		productList = append(productList, formattedProduct)
	}

	data := CategoryProductFormatter{
		CategoryID:   category.ID,
		CategoryName: category.Name,
		Products:     productList,
	}
	return data
}
