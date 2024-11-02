package service

import (
	"api-kasirapp/input"
	"api-kasirapp/models"
	repository2 "api-kasirapp/repository"
)

type DiscountService interface {
	Create(input input.DiscountInput) (models.Discount, error)
	GetByID(id int) (models.Discount, error)
	GetAll() ([]models.Discount, error)
	Update(ID int, input input.DiscountInput) (models.Discount, error)
	Delete(ID int) (models.Discount, error)
}

type discountService struct {
	repository repository2.DiscountRepository
}

func NewDiscountService(repository repository2.DiscountRepository) *discountService {
	return &discountService{repository}
}

func (s *discountService) Create(input input.DiscountInput) (models.Discount, error) {
	discount := models.Discount{
		Name:       input.Name,
		Percentage: input.Percentage,
	}
	newDiscount, err := s.repository.SaveDiscount(discount)
	if err != nil {
		return newDiscount, err
	}
	return newDiscount, nil
}

func (s *discountService) GetByID(id int) (models.Discount, error) {
	discount, err := s.repository.FindDiscountByID(id)
	if err != nil {
		return discount, err
	}
	return discount, nil
}

func (s *discountService) GetAll() ([]models.Discount, error) {
	discounts, err := s.repository.FindDiscounts()
	if err != nil {
		return discounts, err
	}
	return discounts, nil
}

func (s *discountService) Update(ID int, input input.DiscountInput) (models.Discount, error) {
	discount, err := s.repository.FindDiscountByID(ID)
	if err != nil {
		return discount, err
	}

	discount.Name = input.Name
	discount.Percentage = input.Percentage

	updatedDiscount, err := s.repository.UpdateDiscount(ID, discount)
	if err != nil {
		return updatedDiscount, err
	}
	return updatedDiscount, nil
}

func (s *discountService) Delete(ID int) (models.Discount, error) {
	discount, err := s.repository.FindDiscountByID(ID)
	if err != nil {
		return discount, err
	}

	deletedDiscount, err := s.repository.DeleteDiscount(ID)
	if err != nil {
		return deletedDiscount, err
	}
	return deletedDiscount, nil
}
