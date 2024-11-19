package service

import (
	"api-kasirapp/input"
	"api-kasirapp/models"
	"api-kasirapp/repository"
	"errors"
	"strconv"

	"gorm.io/gorm"
)

type OrderServices interface {
	CreateTransactionWithCash(input input.TransactionInput) (models.Transaction, float64, error)
	GetTransactions(ID int) (models.Transaction, error)
}

type orderService struct {
	orderRepository   repository.OrderRepository
	productRepository repository.ProductRepository
}

func NewOrderService(orderRepository repository.OrderRepository, productRepository repository.ProductRepository) *orderService {
	return &orderService{orderRepository, productRepository}
}

func (s *orderService) CreateTransactionWithCash(input input.TransactionInput) (models.Transaction, float64, error) {
	trx := models.Transaction{}
	var details []models.TransactionDetail
	totalCost := 0.0

	for _, productInput := range input.Products {
		product, err := s.productRepository.FindByID(productInput.ProductID)
		if err != nil {
			return trx, 0, err
		}

		if product.Stock < productInput.Qty {
			return trx, 0, errors.New("stock not enough for product ID " + strconv.Itoa(productInput.ProductID))
		}

		// Calculate cost for this product
		productCost := product.SellingPrice * float64(productInput.Qty)
		totalCost += productCost

		// Deduct stock
		product.Stock -= productInput.Qty
		_, err = s.productRepository.Update(product)
		if err != nil {
			return trx, 0, err
		}

		// Add to transaction details
		details = append(details, models.TransactionDetail{
			ProductID: productInput.ProductID,
			Qty:       productInput.Qty,
		})
	}

	// Check if balance is sufficient
	if float64(input.Balance) < totalCost {
		return trx, 0, errors.New("balance not enough")
	}
	cashReturn := float64(input.Balance) - totalCost

	// Save transaction and details
	trx.Amount = totalCost
	trx.Qty = len(input.Products)

	savedTransaction, err := s.orderRepository.Create(trx, details)
	if err != nil {
		return savedTransaction, 0, err
	}

	// Fetch transaction with details
	err = s.orderRepository.GetByIDWithDetails(savedTransaction.ID, &savedTransaction)
	if err != nil {
		return savedTransaction, 0, err
	}

	return savedTransaction, cashReturn, nil

}

func (s *orderService) GetTransactions(ID int) (models.Transaction, error) {
	data, err := s.orderRepository.GetByID(ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, errors.New("transaction not found")
		}
		return data, err
	}

	return data, nil
}
