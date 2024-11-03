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
	GetTransaction(ID int) (*models.Transaction, error)
	// HandleSentEmail(data []byte) error
	// HandleLogging(data []byte) error
	// HandleCallback(notificationPayload map[string]interface{}) error
}

type orderService struct {
	orderRepository repository.OrderRepository
}

func NewOrderService(orderRepository repository.OrderRepository) *orderService {
	return &orderService{orderRepository}
}

func (s *orderService) CreateTransaction(input input.TransactionInput) (*models.Transaction, error) {
	data := models.Transaction{}
	data.ProductID = input.ProductID
	data.Qty = input.Qty
	data.Amount = input.Amount

	newData, err := s.orderRepository.Create(&data)
	if err != nil {
		return newData, err
	}

	return newData, nil
}

func (s *orderService) GetTransaction(ID int) (*models.Transaction, error) {
	data, err := s.orderRepository.GetByID(ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, errors.New("category not found")
		}
		return data, err
	}

	return data, nil
}
