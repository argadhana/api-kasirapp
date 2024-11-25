package service

import (
	"api-kasirapp/input"
	"api-kasirapp/models"
	"api-kasirapp/repository"
	"errors"

	"gorm.io/gorm"
)

type CategoryService interface {
	SaveCategory(input input.CategoryInput) (models.Category, error)
	FindCategoryByID(ID int) (models.Category, error)
	FindCategories() ([]models.Category, error)
	UpdateCategory(ID int, input input.CategoryInput) (models.Category, error)
	DeleteCategory(ID int) (models.Category, error)
	GetCategoryProducts(ID int) ([]models.Product, error)
	GetProductsWithCategoryName(categoryName string) ([]models.Product, error)
	GetCategoryByName(name string) (models.Category, error)
}

type categoryService struct {
	repository repository.CategoryRepository
}

func NewCategoryService(repository repository.CategoryRepository) *categoryService {
	return &categoryService{repository}
}

func (s *categoryService)GetCategoryProducts(ID int) ([]models.Product, error){
	products, err := s.repository.FindCategoryProducts(ID)
	if err != nil {
		return products, err
	}

	return products, nil
}

func (s *categoryService)GetProductsWithCategoryName(categoryName string) ([]models.Product, error){
	products, err := s.repository.FindProductsWithCategoryName(categoryName)
	if err != nil {
		return products, err
	}

	return products, nil
}

func (s *categoryService) SaveCategory(input input.CategoryInput) (models.Category, error) {
	category := models.Category{}

	category.Name = input.Name

	newCategory, err := s.repository.SaveCategory(category)
	if err != nil {
		return newCategory, err
	}

	return newCategory, nil
}

func (s *categoryService) FindCategoryByID(ID int) (models.Category, error) {
	category, err := s.repository.FindCategoryByID(ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {


			return category, errors.New("category not found")
		}
		return category, err
	}

	return category, nil
}

func (s *categoryService) FindCategories() ([]models.Category, error) {
	categories, err := s.repository.FindCategories()
	if err != nil {
		return categories, err
	}

	return categories, nil
}

func (s *categoryService) GetCategoryByName(name string) (models.Category, error) {
	category, err := s.repository.FindCategoryByName(name)
	if err != nil {
		return category, err
	}

	return category, nil

}

func (s *categoryService) UpdateCategory(ID int, input input.CategoryInput) (models.Category, error) {
	category, err := s.repository.FindCategoryByID(ID)
	if err != nil {
		return category, err
	}

	category.Name = input.Name

	updatedCategory, err := s.repository.UpdateCategory(category)
	if err != nil {
		return updatedCategory, err
	}

	return updatedCategory, nil
}

func (s *categoryService) DeleteCategory(ID int) (models.Category, error) {
	category, err := s.repository.FindCategoryByID(ID)
	if err != nil {
		return category, err
	}

	deletedCategory, err := s.repository.DeleteCategory(ID)
	if err != nil {
		return deletedCategory, err
	}

	return deletedCategory, nil
}
