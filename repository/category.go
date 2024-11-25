package repository

import (
	"api-kasirapp/models"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	SaveCategory(category models.Category) (models.Category, error)
	FindCategoryByID(ID int) (models.Category, error)
	FindCategories() ([]models.Category, error)
	FindCategoryByName(name string) (models.Category, error)
	UpdateCategory(category models.Category) (models.Category, error)
	DeleteCategory(ID int) (models.Category, error)
	FindCategoryProducts(ID int) ([]models.Product, error)
	FindProductsWithCategoryName(categoryName string) ([]models.Product, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *categoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository )FindCategoryProducts(ID int) ([]models.Product, error){
	var products []models.Product

	err := r.db.Where("category_id = ?", ID).Find(&products).Error
	if err != nil {
		return products, err
	}
	return products, nil
}

func (r *categoryRepository) SaveCategory(category models.Category) (models.Category, error) {
	err := r.db.Create(&category).Error
	if err != nil {
		return category, err
	}

	return category, nil
}

func (r *categoryRepository) FindCategoryByID(ID int) (models.Category, error) {
	var category models.Category

	err := r.db.Where("id = ?", ID).First(&category).Error
	if err != nil {
		return category, err
	}

	return category, nil
}

func (r *categoryRepository) FindCategories() ([]models.Category, error) {
	var categories []models.Category

	err := r.db.Find(&categories).Error
	if err != nil {
		return categories, err
	}

	return categories, nil
}

func (r *categoryRepository) FindCategoryByName(name string) (models.Category, error) {
	var category models.Category

	err := r.db.Where("name = ?", name).Find(&category).Error
	if err != nil {
		return category, err
	}

	return category, nil
}

func (r *categoryRepository) UpdateCategory(category models.Category) (models.Category, error) {
	err := r.db.Save(&category).Error
	if err != nil {
		return category, err
	}

	return category, nil
}

func (r *categoryRepository) DeleteCategory(ID int) (models.Category, error) {
	var category models.Category

	err := r.db.Where("id = ?", ID).Delete(&category).Error
	if err != nil {
		return category, err
	}
	return category, nil
}

func (r *categoryRepository) FindProductsWithCategoryName(categoryName string) ([]models.Product, error) {
	var category []models.Product

	err := r.db.Preload("Product").Where("name = ?", categoryName).First(&category).Error
	if err != nil {
		return category, err
	}
	return category, nil
}