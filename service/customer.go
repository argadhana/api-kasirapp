package service

import (
	"api-kasirapp/helper"
	"api-kasirapp/input"
	"api-kasirapp/models"
	repository2 "api-kasirapp/repository"
	"errors"

	"gorm.io/gorm"
)

type CustomerService interface {
	CreateCustomer(input input.CustomerInput) (models.Customer, error)
	GetCustomers(limit int, offset int) ([]models.Customer, error)
	GetCustomerByID(ID int) (models.Customer, error)
	UpdateCustomer(ID int, input input.CustomerInput) (models.Customer, error)
	DeleteCustomer(ID int) (models.Customer, error)
	CountCustomers() (int64, error)
}

type customerService struct {
	repository repository2.CustomerRepository
}

func NewCustomerService(repository repository2.CustomerRepository) *customerService {
	return &customerService{repository}
}

func (s *customerService) CreateCustomer(input input.CustomerInput) (models.Customer, error) {
	customer := models.Customer{}

	customer.Name = input.Name
	customer.Address = input.Address
	customer.Phone = input.Phone
	customer.Email = input.Email

	if err := helper.ValidateEmail(customer.Email); err != nil {
		return models.Customer{}, err
	}

	if err := helper.ValidatePhoneNumber(customer.Phone); err != nil {
		return models.Customer{}, err
	}

	newCustomer, err := s.repository.SaveCustomer(customer)
	if err != nil {
		return newCustomer, err
	}

	return newCustomer, nil

}

func (s *customerService) GetCustomers(limit int, offset int) ([]models.Customer, error) {
	customers, err := s.repository.FindCustomers(limit, offset)
	if err != nil {
		return customers, err
	}

	return customers, nil
}

func (s *customerService) GetCustomerByID(ID int) (models.Customer, error) {
	customer, err := s.repository.FindCustomerByID(ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customer, errors.New("category not found")
		}
		return customer, err
	}

	return customer, nil
}

func (s *customerService) UpdateCustomer(ID int, input input.CustomerInput) (models.Customer, error) {
	customer, err := s.repository.FindCustomerByID(ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customer, errors.New("customer not found")
		}
		return customer, err
	}

	customer.Name = input.Name
	customer.Address = input.Address
	customer.Phone = input.Phone
	customer.Email = input.Email

	updatedCustomer, err := s.repository.UpdateCustomer(customer)
	if err != nil {
		return updatedCustomer, err
	}

	return updatedCustomer, nil
}

func (s *customerService) DeleteCustomer(ID int) (models.Customer, error) {
	customer, err := s.repository.FindCustomerByID(ID)
	if err != nil {
		return customer, err
	}

	deletedCustomer, err := s.repository.DeleteCustomer(ID)
	if err != nil {
		return deletedCustomer, err
	}

	return deletedCustomer, nil
}

func (s *customerService) CountCustomers() (int64, error) {
	return s.repository.CountCustomers()
}

