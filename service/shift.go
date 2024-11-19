package service

import (
	"api-kasirapp/input"
	"api-kasirapp/models"
	"api-kasirapp/repository"
)

type ShiftService interface {
	StartShift(input input.ShiftInput) (*models.Shift, error)
	EndShift(input input.ShiftInput) (*models.Shift, error)
}

type shiftService struct {
	shiftRepository repository.ShiftRepository
}

func NewShiftService(shiftRepository repository.ShiftRepository) *shiftService {
	return &shiftService{shiftRepository}
}

func (s *shiftService) StartShift(input input.ShiftInput) (*models.Shift, error) {
	shift := models.Shift{}
	shift.StartBalance = input.StartBalance

	shift, err := s.shiftRepository.Save(shift)
	if err != nil {
		return nil, err
	}

	return &shift, nil
}

func (s *shiftService) EndShift(id int) (*models.Shift, error) {
	shift, err := s.shiftRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	return &shift, nil

}
