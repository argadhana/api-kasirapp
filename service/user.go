package service

import (
	"api-kasirapp/helper"
	"api-kasirapp/input"
	"api-kasirapp/models"
	repository2 "api-kasirapp/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(input input.RegisterUserInput) (models.User, error)
	Login(input input.LoginInput) (models.User, error)
	IsEmailAvailable(input input.CheckEmailInput) (bool, error)
	GetUserByID(ID int) (models.User, error)
	GetAllUsers() ([]models.User, error)
	isActiveUser(ID int) (models.User, error)
}

type userService struct {
	repository repository2.UserRepository
}

func NewService(repository repository2.UserRepository) *userService {
	return &userService{repository}
}

func (s *userService) RegisterUser(input input.RegisterUserInput) (models.User, error) {
	user := models.User{}

	user.Email = input.Email
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}

	user.PasswordHash = string(passwordHash)
	user.Phone = input.Phone

	if err := helper.ValidateEmail(user.Email); err != nil {
		return models.User{}, err
	}

	if err := helper.ValidatePhoneNumber(user.Phone); err != nil {
		return models.User{}, err
	}

	newUser, err := s.repository.Save(user)
	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

func (s *userService) Login(input input.LoginInput) (models.User, error) {
	email := input.Email
	password := input.Password

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("no user found on that email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *userService) IsEmailAvailable(input input.CheckEmailInput) (bool, error) {
	email := input.Email

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return false, err
	}

	if user.ID == 0 {
		return true, nil
	}

	return false, nil
}

func (s *userService) GetUserByID(ID int) (models.User, error) {
	user, err := s.repository.FindByID(ID)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("no user found with that id")
	}

	return user, nil
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	users, err := s.repository.FindAll()
	if err != nil {
		return users, err
	}

	return users, nil
}

func (s *userService) isActiveUser(ID int) (models.User, error) {
	user, err := s.repository.ActivateUser(ID)
	if err != nil {
		return user, err
	}

	return user, nil
}
