// service/stock_service.go
package service

import (
	"api-kasirapp/input"
	"api-kasirapp/models"
	"api-kasirapp/repository"
	"errors"
	"fmt"
)

type StockService interface {
	AddStock(input input.CreateStockInput) (models.Stock, error)
	GetStocks(limit int, offset int) ([]models.Stock, error)
	GetStocksByProductID(productID int) ([]models.Stock, error)
	CountStocks() (int64, error)
	DeleteStock(id int) error
	GetStockByID(id int) (models.Stock, error)
	UpdateStockByID(id int, input input.CreateStockInput) (models.Stock, error)
}

type stockService struct {
	stockrepository   repository.StockRepository
	productRepository repository.ProductRepository
}

func NewStockService(stockRepo repository.StockRepository, productRepo repository.ProductRepository) *stockService {
	return &stockService{
		stockrepository:   stockRepo,
		productRepository: productRepo,
	}
}
func (s *stockService) AddStock(input input.CreateStockInput) (models.Stock, error) {
	stock := models.Stock{
		ProductID:    input.ProductID,
		Quantity:     input.Quantity,
		BasePrice:    input.BasePrice,
		SellingPrice: input.SellingPrice,
		PurchasePrice: input.PurchasePrice,
		Date:         input.Date,
		Description:  input.Description,
	}

	// Fetch the product
	product, err := s.productRepository.FindByID(input.ProductID)
	if err != nil {
		return models.Stock{}, fmt.Errorf("product not found: %w", err)
	}

	// Update product stock
	product.Stock += input.Quantity
	if _, err := s.productRepository.Update(product); err != nil {
		return models.Stock{}, fmt.Errorf("failed to update product stock: %w", err)
	}

	// Create the stock record
	newStock, err := s.stockrepository.Create(stock)
	if err != nil {
		return models.Stock{}, fmt.Errorf("failed to create stock record: %w", err)
	}

	// Ensure the Product is preloaded
	newStock.Product = product

	return newStock, nil
}

func (s *stockService) GetStocks(limit int, offset int) ([]models.Stock, error) {
	return s.stockrepository.FindStocks(limit, offset)
}

func (s *stockService) GetStocksByProductID(productID int) ([]models.Stock, error) {
	stocks, err := s.stockrepository.GetByProductID(productID)
	if err != nil {
		return []models.Stock{}, err
	}
	if len(stocks) == 0 {
		return []models.Stock{}, errors.New("no stock found for the given product")
	}
	return stocks, nil
}

func (s *stockService) CountStocks() (int64, error) {
	return s.stockrepository.CountStocks()
}

func (s *stockService) DeleteStock(id int) error {
	// Check if the stock exists
	stock, err := s.stockrepository.GetByID(id) // Ensure a GetByID method exists
	if err != nil {
		return errors.New("stock not found")
	}

	// Perform the delete operation
	err = s.stockrepository.DeleteByID(stock.ID)
	if err != nil {
		return errors.New("failed to delete stock")
	}

	return nil
}

func (s *stockService) GetStockByID(id int) (models.Stock, error) {
	stock, err := s.stockrepository.GetByID(id)
	if err != nil {
		return models.Stock{}, err
	}
	return stock, nil
}

func (s *stockService) UpdateStockByID(id int, input input.CreateStockInput) (models.Stock, error) {
	// Validate the stock existence
	_, err := s.stockrepository.GetByID(id)
	if err != nil {
		return models.Stock{}, fmt.Errorf("stock not found")
	}

	// Validate product existence (optional, based on business rules)
	_, err = s.productRepository.FindByID(input.ProductID)
	if err != nil {
		return models.Stock{}, fmt.Errorf("product not found")
	}

	// Prepare updated stock data
	updatedStock := models.Stock{
		ProductID:    input.ProductID,
		Quantity:     input.Quantity,
		BasePrice:    input.BasePrice,
		SellingPrice: input.SellingPrice,
		PurchasePrice: input.PurchasePrice,
		Date:         input.Date,
		Description:  input.Description,
	}

	// Call repository to update the stock
	newStock, err := s.stockrepository.UpdateByID(id, updatedStock)
	if err != nil {
		return models.Stock{}, fmt.Errorf("failed to update stock: %w", err)
	}

	return newStock, nil
}
