package repository

import (
	"api-kasirapp/models"
	"errors"

	"gorm.io/gorm"
)

type CustomerRepository interface {
	SaveCustomer(customer models.Customer) (models.Customer, error)
	FindCustomers(limit int, offset int) ([]models.Customer, error)
	FindCustomerByID(ID int) (models.Customer, error)
	UpdateCustomer(customer models.Customer) (models.Customer, error)
	DeleteCustomer(ID int) (models.Customer, error)
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *customerRepository {
	return &customerRepository{db}
}

func (r *customerRepository) SaveCustomer(customer models.Customer) (models.Customer, error) {
	err := r.db.Create(&customer).Error
	if err != nil {
		return customer, err
	}

	return customer, nil
}

func (r *customerRepository) FindCustomers(limit int, offset int) ([]models.Customer, error) {
	var customers []models.Customer

	err := r.db.Limit(limit).Offset(offset).Find(&customers).Error
	if err != nil {
		return customers, err
	}

	return customers, nil
}

func (r *customerRepository) FindCustomerByID(ID int) (models.Customer, error) {
	var customer models.Customer
	err := r.db.Where("id = ?", ID).Find(&customer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customer, errors.New("customer not found")
		}
		return customer, err
	}

	return customer, nil
}

func (r *customerRepository) UpdateCustomer(customer models.Customer) (models.Customer, error) {
	err := r.db.Save(&customer).Error
	if err != nil {
		return customer, err
	}

	return customer, nil
}

func (r *customerRepository) DeleteCustomer(ID int) (models.Customer, error) {
	var customer models.Customer

	err := r.db.Where("id = ?", ID).First(&customer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customer, errors.New("customer not found")
		}
		return customer, err
	}

	return customer, nil
}
