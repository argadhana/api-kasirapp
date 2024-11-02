package service

import (
	"api-kasirapp/helper"
	"api-kasirapp/input"
	"api-kasirapp/models"
	"api-kasirapp/repository"
	"time"

	"golang.org/x/exp/rand"
)

type SupplierService interface {
	CreateSupplier(Input input.SupplierInput) (models.Supplier, error)
	GetByID(ID int) (models.Supplier, error)
	GetByName(name string) (models.Supplier, error)
	GetAll(limit int, offset int) ([]models.Supplier, error)
	Update(ID int, Input input.SupplierInput) (models.Supplier, error)
	Delete(ID int) (models.Supplier, error)
}

type supplierService struct {
	repository repository.SupplierRepository
}

func NewSupplierService(repository repository.SupplierRepository) *supplierService {
	return &supplierService{repository}
}

func (s *supplierService) CreateSupplier(input input.SupplierInput) (models.Supplier, error) {
	supplier := models.Supplier{
		Name:    input.Name,
		Address: input.Address,
		Email:   input.Email,
		Phone:   input.Phone,
	}

	rand.Seed(uint64(time.Now().UnixNano()))
	supplier.Code = rand.Intn(90000) + 10000

	if err := helper.ValidateEmail(supplier.Email); err != nil {
		return models.Supplier{}, err
	}

	if err := helper.ValidatePhoneNumber(supplier.Phone); err != nil {
		return models.Supplier{}, err
	}

	newSupplier, err := s.repository.Save(supplier)
	if err != nil {
		return newSupplier, err
	}

	return newSupplier, nil
}

func (s *supplierService) GetByID(ID int) (models.Supplier, error) {
	supplier, err := s.repository.FindByID(ID)
	if err != nil {
		return supplier, err
	}
	return supplier, nil
}

func (s *supplierService) GetByName(name string) (models.Supplier, error) {
	supplier, err := s.repository.FindByName(name)
	if err != nil {
		return supplier, err
	}
	return supplier, nil
}

func (s *supplierService) GetAll(limit int, offset int) ([]models.Supplier, error) {
	suppliers, err := s.repository.FindAll(limit, offset)
	if err != nil {
		return suppliers, err
	}
	return suppliers, nil
}

func (s *supplierService) Update(ID int, input input.SupplierInput) (models.Supplier, error) {
	supplier, err := s.repository.FindByID(ID)
	if err != nil {
		return supplier, err
	}

	supplier.Name = input.Name
	supplier.Address = input.Address
	supplier.Email = input.Email
	supplier.Phone = input.Phone

	updatedSupplier, err := s.repository.Update(ID, supplier)
	if err != nil {
		return updatedSupplier, err
	}
	return updatedSupplier, nil
}

func (s *supplierService) Delete(ID int) (models.Supplier, error) {
	supplier, err := s.repository.FindByID(ID)
	if err != nil {
		return supplier, err
	}
	deletedSupplier, err := s.repository.Delete(ID)
	if err != nil {
		return deletedSupplier, err
	}
	return deletedSupplier, nil
}
