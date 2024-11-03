// repository/stock_repository.go
package repository

import (
	"api-kasirapp/models"

	"gorm.io/gorm"
)

type StockRepository interface {
	Create(stock models.Stock) (models.Stock, error)
	GetAll() ([]models.Stock, error)
	GetByProductID(productID int) ([]models.Stock, error)
}

type stockRepository struct {
	db *gorm.DB
}

func NewStockRepository(db *gorm.DB) *stockRepository {
	return &stockRepository{db}
}

func (r *stockRepository) Create(stock models.Stock) (models.Stock, error) {
	err := r.db.Create(&stock).Error
	return stock, err
}

func (r *stockRepository) GetAll() ([]models.Stock, error) {
	var stocks []models.Stock
	err := r.db.Preload("Product").Find(&stocks).Error
	return stocks, err
}

func (r *stockRepository) GetByProductID(productID int) ([]models.Stock, error) {
	var stocks []models.Stock
	err := r.db.Where("product_id = ?", productID).Preload("Product").Find(&stocks).Error
	return stocks, err
}
