package repository

import (
	"api-kasirapp/models"
	"errors"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Save(product models.Product) (models.Product, error)
	FindByID(ID int) (models.Product, error)
	FindByName(name string) (models.Product, error)
	FindAll() ([]models.Product, error)
	FindByCategoryID(categoryID int) ([]models.Product, error)
	Update(product models.Product) (models.Product, error)
	Delete(ID int) (models.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *productRepository {
	return &productRepository{db}
}

func (r *productRepository) FindByCategoryID(categoryID int) ([]models.Product, error) {
	var products []models.Product

	err := r.db.Where("category_id = ?", categoryID).Find(&products).Error
	if err != nil {
		return products, err
	}
	return products, nil
}

func (r *productRepository) Save(product models.Product) (models.Product, error) {
	var existingProduct models.Product

	// Check if the product code already exists
	if err := r.db.Where("code_product = ?", product.CodeProduct).First(&existingProduct).Error; err == nil {
		return product, errors.New("product code already exists") // Return error if product code exists
	}

	// Insert new product
	if err := r.db.Create(&product).Error; err != nil {
		return product, err
	}

	return product, nil
}

func (r *productRepository) FindByID(productID int) (models.Product, error) {
	var product models.Product
	err := r.db.First(&product, "id = ?", productID).Error
	if err != nil {
		return product, err
	}
	return product, nil
}

func (r *productRepository) FindByName(name string) (models.Product, error) {
	var product models.Product

	err := r.db.Where("name = ?", name).Find(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return product, errors.New("product not found")
		}
		return product, err
	}
	return product, nil
}

func (r *productRepository) FindAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}
func (r *productRepository) Update(product models.Product) (models.Product, error) {
	if err := r.db.Save(&product).Error; err != nil {
		return product, err
	}
	return product, nil
}

func (r *productRepository) Delete(ID int) (models.Product, error) {
	var product models.Product

	err := r.db.Where("id = ?", ID).Delete(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return product, errors.New("product not found")
		}
		return product, err
	}

	return product, nil
}
