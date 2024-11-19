package service

import (
	"api-kasirapp/input"
	"api-kasirapp/models"
	"api-kasirapp/repository"
	"errors"
	"time"
)

type ShiftService interface {
	StartShift(input input.ShiftInput) (*models.Shift, error)
	EndShift(ID int) (*models.Shift, error)
}

type shiftService struct {
	shiftRepository repository.ShiftRepository
	orderRepository repository.OrderRepository
}

func NewShiftService(shiftRepository repository.ShiftRepository, orderRepository repository.OrderRepository) ShiftService {
	return &shiftService{
		shiftRepository: shiftRepository,
		orderRepository: orderRepository,
	}
}

func (s *shiftService) StartShift(input input.ShiftInput) (*models.Shift, error) {
	shift := models.Shift{}
	shift.StartBalance = input.StartBalance
	shift.StartTime = time.Now()

	shift, err := s.shiftRepository.Save(shift)
	if err != nil {
		return nil, err
	}

	return &shift, nil
}

func (s *shiftService) EndShift(ID int) (*models.Shift, error) {
	shift, err := s.shiftRepository.FindByID(ID)
	if err != nil {
		return nil, err
	}

	if shift.Status != "berjalan" {
		return nil, errors.New("shift is not running")
	}

	totalSales, err := s.orderRepository.GetTotalSalesByShiftID(ID)
	if err != nil {
		return nil, err
	}
	endTime := time.Now()

	shift.Status = "selesai"
	shift.TotalSales = totalSales
	shift.EndTime = &endTime

	shift, err = s.shiftRepository.Update(ID, shift)
	if err != nil {
		return nil, err
	}

	return &shift, nil
}
