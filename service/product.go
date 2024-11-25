package service

import (
	"api-kasirapp/input"
	"api-kasirapp/models"
	"api-kasirapp/repository"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
	"time"
)

type ProductService interface {
	CreateProduct(input input.ProductInput) (models.Product, error)
	FindProductByID(ID int) (models.Product, error)
	FindByName(name string) (models.Product, error)
	FindAll() ([]models.Product, error)
	UpdateProduct(ID int, input input.ProductInput) (models.Product, error)
	DeleteProduct(ID int) (models.Product, error)
	ExportProductsToXLS() (*excelize.File, error)
	ImportProductsFromXLS(filePath string) ([]models.Product, error)
	SaveProductImage(ID int, fileLocation string) (models.Product, error)
}

type productService struct {
	productRepository  repository.ProductRepository
	categoryRepository repository.CategoryRepository
}

func NewProductService(productRepository repository.ProductRepository, categoryRepository repository.CategoryRepository) *productService {
	return &productService{productRepository, categoryRepository}
}

func (s *productService) CreateProduct(input input.ProductInput) (models.Product, error) {
	product := models.Product{}

	product.Name = input.Name
	product.ProductType = input.ProductType
	product.BasePrice = input.BasePrice
	product.SellingPrice = input.SellingPrice
	product.Stock = input.Stock
	product.CodeProduct = input.CodeProduct
	product.CategoryID = input.CategoryID
	product.MinimumStock = input.MinimumStock
	product.Shelf = input.Shelf
	product.Weight = input.Weight
	product.Discount = input.Discount
	product.Information = input.Information

	newProduct, err := s.productRepository.Save(product)
	if err != nil {
		return newProduct, err
	}

	return newProduct, nil
}

func (s *productService) FindProductByID(ID int) (models.Product, error) {
	product, err := s.productRepository.FindByID(ID)
	if err != nil {
		return product, err
	}

	return product, nil
}

func (s *productService) FindByName(name string) (models.Product, error) {
	product, err := s.productRepository.FindByName(name)
	if err != nil {
		return product, err
	}

	return product, nil
}

func (s *productService) FindAll() ([]models.Product, error) {
	products, err := s.productRepository.FindAll()
	if err != nil {
		return products, err
	}

	return products, nil
}

func (s *productService) UpdateProduct(ID int, input input.ProductInput) (models.Product, error) {
	product, err := s.productRepository.FindByID(ID)
	if err != nil {
		return product, err
	}

	product.Name = input.Name
	product.ProductType = input.ProductType
	product.BasePrice = input.BasePrice
	product.SellingPrice = input.SellingPrice
	product.Stock = input.Stock
	product.CodeProduct = input.CodeProduct
	product.CategoryID = input.CategoryID
	product.MinimumStock = input.MinimumStock
	product.Shelf = input.Shelf
	product.Weight = input.Weight
	product.Discount = input.Discount
	product.Information = input.Information

	updatedProduct, err := s.productRepository.Update(product)
	if err != nil {
		return updatedProduct, err
	}

	return updatedProduct, nil
}

func (s *productService) DeleteProduct(ID int) (models.Product, error) {
	product, err := s.productRepository.FindByID(ID)
	if err != nil {
		return product, err
	}

	deletedProduct, err := s.productRepository.Delete(ID)
	if err != nil {
		return deletedProduct, err
	}

	return deletedProduct, nil
}

func (s *productService) ExportProductsToXLS() (*excelize.File, error) {
	// Fetch products from the database
	products, err := s.productRepository.FindAll()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Number of products retrieved: %d\n", len(products))

	if len(products) == 0 {
		return nil, errors.New("no products found to export")
	}

	// Create a new Excel file
	f := excelize.NewFile()
	sheet := "Products"
	index, err := f.NewSheet(sheet)
	if err != nil {
		return nil, err
	}

	// Write headers
	headers := []string{"ID", "Name", "Product Type", "Base Price", "Selling Price", "Stock", "Code Product", "Category ID", "Minimum Stock", "Shelf", "Weight", "Discount", "Information", "Created At", "Updated At"}
	for i, header := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheet, cell, header)
	}

	// Write product data
	for i, product := range products {
		row := i + 2
		fmt.Printf("Writing product %d to row %d\n", product.ID, row) // Debugging output

		f.SetCellValue(sheet, "A"+strconv.Itoa(row), product.ID)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), product.Name)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), product.ProductType)
		f.SetCellValue(sheet, "D"+strconv.Itoa(row), product.BasePrice)
		f.SetCellValue(sheet, "E"+strconv.Itoa(row), product.SellingPrice)
		f.SetCellValue(sheet, "F"+strconv.Itoa(row), product.Stock)
		f.SetCellValue(sheet, "G"+strconv.Itoa(row), product.CodeProduct)
		f.SetCellValue(sheet, "H"+strconv.Itoa(row), product.CategoryID)
		f.SetCellValue(sheet, "I"+strconv.Itoa(row), product.MinimumStock)
		f.SetCellValue(sheet, "J"+strconv.Itoa(row), product.Shelf)
		f.SetCellValue(sheet, "K"+strconv.Itoa(row), product.Weight)
		f.SetCellValue(sheet, "L"+strconv.Itoa(row), product.Discount)
		f.SetCellValue(sheet, "M"+strconv.Itoa(row), product.Information)
		f.SetCellValue(sheet, "N"+strconv.Itoa(row), product.CreatedAt.Format(time.RFC3339))
		f.SetCellValue(sheet, "O"+strconv.Itoa(row), product.UpdatedAt.Format(time.RFC3339))
	}

	// Set active sheet and delete the default sheet
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	return f, nil
}

func (s *productService) ImportProductsFromXLS(filePath string) ([]models.Product, error) {
	// Open the Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read data from the first sheet
	sheet := "Products"
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}

	var importedProducts []models.Product

	// Skip the header row
	for i, row := range rows {
		if i == 0 {
			continue
		}

		// Parse each column into product fields
		basePrice, _ := strconv.ParseFloat(row[3], 64)
		sellingPrice, _ := strconv.ParseFloat(row[4], 64)
		stock, _ := strconv.Atoi(row[5])
		categoryID, _ := strconv.Atoi(row[7])
		minimumStock, _ := strconv.Atoi(row[8])
		weight, _ := strconv.Atoi(row[10])
		discount, _ := strconv.Atoi(row[11])

		product := models.Product{
			Name:         row[1],
			ProductType:  row[2],
			BasePrice:    basePrice,
			SellingPrice: sellingPrice,
			Stock:        stock,
			CodeProduct:  row[6],
			CategoryID:   categoryID,
			MinimumStock: minimumStock,
			Shelf:        row[9],
			Weight:       weight,
			Discount:     discount,
			Information:  row[12],
		}

		// Insert product into the database
		savedProduct, err := s.productRepository.Save(product)
		if err != nil {
			return nil, err
		}

		// Add the saved product to the list
		importedProducts = append(importedProducts, savedProduct)
	}

	return importedProducts, nil
}

func (s *productService) SaveProductImage(productID int, filePath string) (models.Product, error) {
	// Find the product by ID
	product, err := s.productRepository.FindByID(productID)
	if err != nil {
		return product, fmt.Errorf("product not found: %w", err)
	}

	// Update the product's image URL
	product.ProductFileName = filePath

	// Save the updated product in the database
	updatedProduct, err := s.productRepository.Update(product)
	if err != nil {
		return updatedProduct, fmt.Errorf("failed to update product: %w", err)
	}

	return updatedProduct, nil
}
