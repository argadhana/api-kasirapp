package repository

import (
	"api-kasirapp/models"
	"gorm.io/gorm"
)

type SupplierRepository interface {
	Save(supplier models.Supplier) (models.Supplier, error)
	FindByID(ID int) (models.Supplier, error)
	FindByName(name string) (models.Supplier, error)
	FindAll(limit int, offset int) ([]models.Supplier, error)
	Update(ID int, supplier models.Supplier) (models.Supplier, error)
	Delete(ID int) (models.Supplier, error)
}

type supplierRepository struct {
	db *gorm.DB
}

func NewSupplierRepository(db *gorm.DB) *supplierRepository {
	return &supplierRepository{db}
}

func (r *supplierRepository) Save(supplier models.Supplier) (models.Supplier, error) {
	if err := r.db.Create(&supplier).Error; err != nil {
		return supplier, err
	}

	return supplier, nil
}

func (r *supplierRepository) FindByID(ID int) (models.Supplier, error) {
	var supplier models.Supplier

	err := r.db.Where("id = ?", ID).Find(&supplier).Error
	if err != nil {
		return supplier, err
	}
	return supplier, nil
}

func (r *supplierRepository) FindByName(name string) (models.Supplier, error) {
	var supplier models.Supplier

	err := r.db.Where("name = ?", name).Find(&supplier).Error
	if err != nil {
		return supplier, err
	}
	return supplier, nil
}

func (r *supplierRepository) FindAll(limit int, offset int) ([]models.Supplier, error) {
	var suppliers []models.Supplier

	err := r.db.Limit(limit).Offset(offset).Find(&suppliers).Error
	if err != nil {
		return suppliers, err
	}

	return suppliers, nil
}

func (r *supplierRepository) Update(ID int, supplier models.Supplier) (models.Supplier, error) {
	if err := r.db.Model(&models.Supplier{}).Where("id = ?", ID).Updates(&supplier).Error; err != nil {
		return supplier, err
	}
	return supplier, nil
}

func (r *supplierRepository) Delete(ID int) (models.Supplier, error) {
	var supplier models.Supplier

	err := r.db.Where("id = ?", ID).Delete(&supplier).Error
	if err != nil {
		return supplier, err
	}

	return supplier, nil
}
