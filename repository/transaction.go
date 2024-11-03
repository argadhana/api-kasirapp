package repository

import (
	"api-kasirapp/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(data *models.Transaction) (*models.Transaction, error)
	GetByID(ID int) (*models.Transaction, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *orderRepository {
	return &orderRepository{db}
}

func (r *orderRepository) Create(data *models.Transaction) (*models.Transaction, error){
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(data).Error; err != nil {
		tx.Rollback()
		return data, err
	}

	tx.Commit()

	return data, nil
}

func (r *orderRepository) GetByID(ID int) (*models.Transaction, error) {
	var data models.Transaction

	if err := r.db.Debug().Preload("Product").Model(&models.Transaction{}).Where("id = ?", ID).Find(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}
