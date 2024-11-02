package service

import (
	"api-kasirapp/input"
	"api-kasirapp/models"
	"api-kasirapp/repository"
)

type ProductService interface {
	CreateProduct(input input.ProductInput) (models.Product, error)
	FindProductByID(ID int) (models.Product, error)
	FindByName(name string) (models.Product, error)
	FindAll() ([]models.Product, error)
	UpdateProduct(ID int, input input.ProductInput) (models.Product, error)
	DeleteProduct(ID int) (models.Product, error)
}

type productService struct {
	productRepository repository.ProductRepository
}

func NewProductService(productRepository repository.ProductRepository) *productService {
	return &productService{productRepository}
}

func (s *productService) CreateProduct(input input.ProductInput) (models.Product, error) {
	product := models.Product{}

	product.Name = input.Name
	product.ProductType = input.ProductType
	product.BasePrice = input.BasePrice
	product.SellingPrice = input.SellingPrice
	product.Stock = input.Stock
	product.CodeProduct = input.CodeProduct
	product.CategoryID = input.CategoryID
	product.MinimumStock = input.MinimumStock
	product.Shelf = input.Shelf
	product.Weight = input.Weight
	product.Discount = input.Discount
	product.Information = input.Information

	newProduct, err := s.productRepository.Save(product)
	if err != nil {
		return newProduct, err
	}

	return newProduct, nil
}

func (s *productService) FindProductByID(ID int) (models.Product, error) {
	product, err := s.productRepository.FindByID(ID)
	if err != nil {
		return product, err
	}

	return product, nil
}

func (s *productService) FindByName(name string) (models.Product, error) {
	product, err := s.productRepository.FindByName(name)
	if err != nil {
		return product, err
	}

	return product, nil
}

func (s *productService) FindAll() ([]models.Product, error) {
	products, err := s.productRepository.FindAll()
	if err != nil {
		return products, err
	}

	return products, nil
}

func (s *productService) UpdateProduct(ID int, input input.ProductInput) (models.Product, error) {
	product, err := s.productRepository.FindByID(ID)
	if err != nil {
		return product, err
	}

	product.Name = input.Name
	product.ProductType = input.ProductType
	product.BasePrice = input.BasePrice
	product.SellingPrice = input.SellingPrice
	product.Stock = input.Stock
	product.CodeProduct = input.CodeProduct
	product.CategoryID = input.CategoryID
	product.MinimumStock = input.MinimumStock
	product.Shelf = input.Shelf
	product.Weight = input.Weight
	product.Discount = input.Discount
	product.Information = input.Information

	updatedProduct, err := s.productRepository.Update(ID, product)
	if err != nil {
		return updatedProduct, err
	}

	return updatedProduct, nil
}

func (s *productService) DeleteProduct(ID int) (models.Product, error) {
	product, err := s.productRepository.FindByID(ID)
	if err != nil {
		return product, err
	}

	deletedProduct, err := s.productRepository.Delete(ID)
	if err != nil {
		return deletedProduct, err
	}

	return deletedProduct, nil
}
