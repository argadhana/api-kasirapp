// repository/stock_repository.go
package repository

import (
	"api-kasirapp/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type StockRepository interface {
	Create(stock models.Stock) (models.Stock, error)
	FindStocks(limit int, offset int) ([]models.Stock, error)
	GetByProductID(productID int) ([]models.Stock, error)
	CountStocks() (int64, error)
	DeleteByID(id int) error
	GetByID(id int) (models.Stock, error)
	UpdateByID(id int, stock models.Stock) (models.Stock, error)
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

func (r *stockRepository) DeleteByID(id int) error {
	err := r.db.Delete(&models.Stock{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *stockRepository) GetByID(id int) (models.Stock, error) {
	var stock models.Stock

	// Use GORM to find the stock by ID and preload the associated Product
	err := r.db.Preload("Product").First(&stock, id).Error
	if err != nil {
		// Return an error if the stock is not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return stock, fmt.Errorf("stock not found")
		}
		return stock, err
	}

	return stock, nil
}

func (r *stockRepository) UpdateByID(id int, stock models.Stock) (models.Stock, error) {
	// Find the existing stock
	var existingStock models.Stock
	err := r.db.First(&existingStock, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return existingStock, fmt.Errorf("stock not found")
		}
		return existingStock, err
	}

	// Update the stock
	err = r.db.Model(&existingStock).Updates(stock).Error
	if err != nil {
		return existingStock, err
	}

	// Preload the associated product
	err = r.db.Preload("Product").First(&existingStock, id).Error
	if err != nil {
		return existingStock, err
	}

	return existingStock, nil
}
