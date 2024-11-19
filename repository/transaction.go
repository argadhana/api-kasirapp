package repository

import (
	"api-kasirapp/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(data models.Transaction, details []models.TransactionDetail) (models.Transaction, error)
	GetByIDWithDetails(id int, transaction *models.Transaction) error
	GetByID(ID int) (models.Transaction, error)
	GetTotalSalesByShiftID(ID int) (float64, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *orderRepository {
	return &orderRepository{db}
}

func (r *orderRepository) Create(data models.Transaction, details []models.TransactionDetail) (models.Transaction, error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return data, err
	}

	for _, detail := range details {
		detail.TransactionID = data.ID
		if err := tx.Create(&detail).Error; err != nil {
			tx.Rollback()
			return data, err
		}
	}

	tx.Commit()

	return data, nil
}

func (r *orderRepository) GetByIDWithDetails(id int, transaction *models.Transaction) error {
	return r.db.Preload("Details.Product").First(transaction, id).Error
}


func (r *orderRepository) GetByID(ID int) (models.Transaction, error) {
	var data models.Transaction

	if err := r.db.Debug().Preload("Product").First(&data, ID).Error; err != nil {
		return data, err
	}

	return data, nil
}

func (r *orderRepository) GetTotalSalesByShiftID(ID int) (float64, error) {
	var total float64

	if err := r.db.Debug().Model(&models.Transaction{}).Where("shift_id = ?", ID).Select("COALESCE(SUM(total), 0)").Scan(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}
