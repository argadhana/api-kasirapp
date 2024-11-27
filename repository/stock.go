// repository/stock_repository.go
package repository

import (
	"api-kasirapp/models"
	"fmt"

	"gorm.io/gorm"
)

type StockRepository interface {
	Create(stock models.Stock) (models.Stock, error)
	FindStocks(limit int, offset int) ([]models.Stock, error)
	GetByProductID(productID int) ([]models.Stock, error)
	CountStocks() (int64, error)
}

type stockRepository struct {
	db *gorm.DB
}

func NewStockRepository(db *gorm.DB) *stockRepository {
	return &stockRepository{db}
}

func (r *stockRepository) Create(stock models.Stock) (models.Stock, error) {
	// Create the stock record
	err := r.db.Create(&stock).Error
	if err != nil {
		return stock, err
	}

	// Preload and debug
	err = r.db.Preload("Product").First(&stock, stock.ID).Error
	if err != nil {
		return stock, err
	}

	fmt.Printf("Debug: Stock with Product: %+v\n", stock)
	return stock, nil
}

func (r *stockRepository) FindStocks(limit int, offset int) ([]models.Stock, error) {
	var stocks []models.Stock
	err := r.db.Preload("Product").Limit(limit).Offset(offset).Find(&stocks).Error
	if err != nil {
		return nil, err
	}
	return stocks, nil
}

func (r *stockRepository) GetByProductID(productID int) ([]models.Stock, error) {
	var stocks []models.Stock
	err := r.db.Where("product_id = ?", productID).Preload("Product").Find(&stocks).Error
	if err != nil {
		return nil, err
	}
	return stocks, nil
}

func (r *stockRepository) CountStocks() (int64, error) {
	var total int64
	err := r.db.Model(&models.Stock{}).Count(&total).Error
	return total, err
}
