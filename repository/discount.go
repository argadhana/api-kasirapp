package repository

import (
	"api-kasirapp/models"
	"gorm.io/gorm"
)

type DiscountRepository interface {
	SaveDiscount(discount models.Discount) (models.Discount, error)
	FindDiscountByID(id int) (models.Discount, error)
	FindDiscounts() ([]models.Discount, error)
	UpdateDiscount(ID int, discount models.Discount) (models.Discount, error)
	DeleteDiscount(ID int) (models.Discount, error)
}

type discountRepository struct {
	db *gorm.DB
}

func NewDiscountRepository(db *gorm.DB) *discountRepository {
	return &discountRepository{db}
}

func (r *discountRepository) SaveDiscount(discount models.Discount) (models.Discount, error) {
	err := r.db.Create(&discount).Error
	if err != nil {
		return discount, err
	}
	return discount, nil
}

func (r *discountRepository) FindDiscountByID(id int) (models.Discount, error) {
	var discount models.Discount
	err := r.db.Where("id = ?", id).First(&discount).Error
	if err != nil {
		return discount, err
	}
	return discount, nil
}

func (r *discountRepository) FindDiscounts() ([]models.Discount, error) {
	var discounts []models.Discount
	err := r.db.Find(&discounts).Error
	if err != nil {
		return discounts, err
	}
	return discounts, nil
}

func (r *discountRepository) UpdateDiscount(ID int, input models.Discount) (models.Discount, error) {
	var discount models.Discount
	if err := r.db.Where("id = ?", ID).First(&discount).Error; err != nil {
		return discount, err
	}
	discount.Name = input.Name
	discount.Percentage = input.Percentage
	err := r.db.Save(&discount).Error
	if err != nil {
		return discount, err
	}

	return discount, nil
}

func (r *discountRepository) DeleteDiscount(ID int) (models.Discount, error) {
	var discount models.Discount

	err := r.db.Where("id = ?", ID).Delete(&discount).Error
	if err != nil {
		return discount, err
	}

	return discount, nil
}
