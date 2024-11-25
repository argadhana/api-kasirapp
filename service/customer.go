package service

import (
	"api-kasirapp/helper"
	"api-kasirapp/input"
	"api-kasirapp/models"
	repository2 "api-kasirapp/repository"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type CustomerService interface {
	CreateCustomer(input input.CustomerInput) (models.Customer, error)
	GetCustomers(limit int, offset int) ([]models.Customer, error)
	GetCustomerByID(ID int) (models.Customer, error)
	UpdateCustomer(ID int, input input.CustomerInput) (models.Customer, error)
	DeleteCustomer(ID int) (models.Customer, error)
	CountCustomers() (int64, error)
	ExportCustomersToXLS() (*excelize.File, error)
	ImportCustomersFromXLS(filePath string) ([]models.Customer, error)
}

type customerService struct {
	repository repository2.CustomerRepository
}

func NewCustomerService(repository repository2.CustomerRepository) *customerService {
	return &customerService{repository}
}

func (s *customerService) CreateCustomer(input input.CustomerInput) (models.Customer, error) {
	customer := models.Customer{}

	customer.Name = input.Name
	customer.Address = input.Address
	customer.Phone = input.Phone
	customer.Email = input.Email

	if err := helper.ValidateEmail(customer.Email); err != nil {
		return models.Customer{}, err
	}

	if err := helper.ValidatePhoneNumber(customer.Phone); err != nil {
		return models.Customer{}, err
	}

	newCustomer, err := s.repository.SaveCustomer(customer)
	if err != nil {
		return newCustomer, err
	}

	return newCustomer, nil

}

func (s *customerService) GetCustomers(limit int, offset int) ([]models.Customer, error) {
	customers, err := s.repository.FindCustomers(limit, offset)
	if err != nil {
		return customers, err
	}

	return customers, nil
}

func (s *customerService) GetCustomerByID(ID int) (models.Customer, error) {
	customer, err := s.repository.FindCustomerByID(ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customer, errors.New("category not found")
		}
		return customer, err
	}

	return customer, nil
}

func (s *customerService) UpdateCustomer(ID int, input input.CustomerInput) (models.Customer, error) {
	customer, err := s.repository.FindCustomerByID(ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customer, errors.New("customer not found")
		}
		return customer, err
	}

	customer.Name = input.Name
	customer.Address = input.Address
	customer.Phone = input.Phone
	customer.Email = input.Email

	updatedCustomer, err := s.repository.UpdateCustomer(customer)
	if err != nil {
		return updatedCustomer, err
	}

	return updatedCustomer, nil
}

func (s *customerService) DeleteCustomer(ID int) (models.Customer, error) {
	customer, err := s.repository.FindCustomerByID(ID)
	if err != nil {
		return customer, err
	}

	deletedCustomer, err := s.repository.DeleteCustomer(ID)
	if err != nil {
		return deletedCustomer, err
	}

	return deletedCustomer, nil
}

func (s *customerService) CountCustomers() (int64, error) {
	return s.repository.CountCustomers()
}

func (s *customerService) ExportCustomersToXLS() (*excelize.File, error) {
	// Fetch customers from the database
	customers, err := s.repository.FindCustomers(6, 0)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Number of customers retrieved: %d\n", len(customers))

	if len(customers) == 0 {
		return nil, errors.New("no customers found to export")
	}

	// Create a new Excel file
	f := excelize.NewFile()
	sheet := "Customers"
	index, err := f.NewSheet(sheet)
	if err != nil {
		return nil, err
	}

	// Write headers
	headers := []string{"ID", "Name", "Address", "Phone", "Email", "Created At", "Updated At"}
	for i, header := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheet, cell, header)
	}

	// Write customer data
	for i, customer := range customers {
		row := i + 2
		fmt.Printf("Writing customer %d to row %d\n", customer.ID, row) // Debugging output

		f.SetCellValue(sheet, "A"+strconv.Itoa(row), customer.ID)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), customer.Name)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), customer.Address)
		f.SetCellValue(sheet, "D"+strconv.Itoa(row), customer.Phone)
		f.SetCellValue(sheet, "E"+strconv.Itoa(row), customer.Email)
		f.SetCellValue(sheet, "F"+strconv.Itoa(row), customer.CreatedAt.Format(time.RFC3339))
		f.SetCellValue(sheet, "G"+strconv.Itoa(row), customer.UpdatedAt.Format(time.RFC3339))
	}

	// Set active sheet and delete the default sheet
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	return f, nil
}

func (s *customerService) ImportCustomersFromXLS(filePath string) ([]models.Customer, error) {
	// Open the Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read data from the first sheet
	sheet := "Customers"
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}

	var importedCustomers []models.Customer

	// Skip the header row
	for i, row := range rows {
		if i == 0 {
			continue
		}

		// Parse each column into customer fields
		customer := models.Customer{
			Name:    row[1],
			Address: row[2],
			Phone:   row[3],
			Email:   row[4],
		}

		// Insert customer into the database
		savedCustomer, err := s.repository.SaveCustomer(customer)
		if err != nil {
			return nil, err
		}

		// Add the saved customer to the list
		importedCustomers = append(importedCustomers, savedCustomer)
	}

	return importedCustomers, nil
}
