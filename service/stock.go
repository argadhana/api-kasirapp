// service/stock_service.go
package service

import (
	"api-kasirapp/input"
	"api-kasirapp/models"
	"api-kasirapp/repository"
	"errors"
)

type StockService interface {
	AddStock(input input.CreateStockInput) (models.Stock, error)
	GetStocks() ([]models.Stock, error)
	GetStocksByProductID(productID int) ([]models.Stock, error)
}

type stockService struct {
	repository repository.StockRepository
}

func NewStockService(repository repository.StockRepository) *stockService {
	return &stockService{repository}
}

func (s *stockService) AddStock(input input.CreateStockInput) (models.Stock, error) {
	stock := models.Stock{
		ProductID:    input.ProductID,
		Quantity:     input.Quantity,
		BasePrice:    input.BasePrice,
		SellingPrice: input.SellingPrice,
		Date:         input.Date,
		Description:  input.Description,
	}

	newStock, err := s.repository.Create(stock)
	if err != nil {
		return models.Stock{}, err
	}

	return newStock, nil
}

func (s *stockService) GetStocks() ([]models.Stock, error) {
	return s.repository.GetAll()
}

func (s *stockService) GetStocksByProductID(productID int) ([]models.Stock, error) {
	stocks, err := s.repository.GetByProductID(productID)
	if err != nil {
		return []models.Stock{}, err
	}
	if len(stocks) == 0 {
		return []models.Stock{}, errors.New("no stock found for the given product")
	}
	return stocks, nil
}

