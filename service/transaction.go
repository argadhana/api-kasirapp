package service

import (
	"api-kasirapp/input"
	"api-kasirapp/models"
	"api-kasirapp/repository"
	"errors"

	"gorm.io/gorm"
)

type OrderServices interface {
	CreateTransaction(input input.TransactionInput) (*models.Transaction, error)
	GetTransactions(ID int) (*models.Transaction, error)
	GetTransactionWithProducts(ID int) (*[]models.Product, error)
	// HandleSentEmail(data []byte) error
	// HandleLogging(data []byte) error
	// HandleCallback(notificationPayload map[string]interface{}) error
}

type orderService struct {
	orderRepository   repository.OrderRepository
	productRepository repository.ProductRepository
}

func NewOrderService(orderRepository repository.OrderRepository, productRepository repository.ProductRepository) *orderService {
	return &orderService{orderRepository, productRepository}
}

func (s *orderService) CreateTransaction(input input.TransactionInput) (*models.Transaction, error) {
	trx := &models.Transaction{}

	trx.ProductID = input.ProductID
	trx.Qty = input.Qty
	trx.Amount = input.Amount

	product, err := s.productRepository.FindByID(input.ProductID)
	if err != nil {
		return nil, err
	}

	if product.Stock < input.Qty {
		return nil, errors.New("stock not enough")
	}

	totalCost := product.SellingPrice * float64(input.Qty)
	if float64(input.Balance) < totalCost {
		return nil, errors.New("balance not enough")
	}

	data, err := s.orderRepository.Create(trx)
	if err != nil {
		return nil, err
	}

	return data, nil

}

func (s *orderService) GetTransactions(ID int) (*models.Transaction, error) {
	data, err := s.orderRepository.GetByID(ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, errors.New("transaction not found")
		}
		return data, err
	}

	return data, nil
}

func (s *orderService) GetTransactionWithProducts(ID int) (*[]models.Product, error) {
	products, err := s.orderRepository.GetTransactionWithProducts(ID)
	if err != nil {
		return nil, err
	}

	return &products, nil
}
