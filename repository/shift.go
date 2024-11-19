package repository

import (
	"api-kasirapp/models"

	"gorm.io/gorm"
)

type ShiftRepository interface {
	Save(shift models.Shift) (models.Shift, error)
	FindByID(ID int) (models.Shift, error)
	FindAll() ([]models.Shift, error)
	Update(ID int, shift models.Shift) (models.Shift, error)
}

type shiftRepository struct {
	db *gorm.DB
}

func NewShiftRepository(db *gorm.DB) *shiftRepository {
	return &shiftRepository{db}
}

func (r *shiftRepository) Save(shift models.Shift) (models.Shift, error) {
	if err := r.db.Create(&shift).Error; err != nil {
		return shift, err
	}

	return shift, nil
}

func (r *shiftRepository) FindByID(ID int) (models.Shift, error) {
	var shift models.Shift
	if err := r.db.First(&shift, ID).Error; err != nil {
		return shift, err
	}

	return shift, nil
}

func (r *shiftRepository) FindAll() ([]models.Shift, error) {
	var shifts []models.Shift
	if err := r.db.Find(&shifts).Error; err != nil {
		return nil, err
	}

	return shifts, nil
}

func (r *shiftRepository) Update(ID int, shift models.Shift) (models.Shift, error) {
	if err := r.db.Save(&shift).Error; err != nil {
		return shift, err
	}

	return shift, nil
}
