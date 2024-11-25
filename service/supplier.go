package service

import (
	"api-kasirapp/helper"
	"api-kasirapp/input"
	"api-kasirapp/models"
	"api-kasirapp/repository"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
	"time"

	"golang.org/x/exp/rand"
)

type SupplierService interface {
	CreateSupplier(Input input.SupplierInput) (models.Supplier, error)
	GetByID(ID int) (models.Supplier, error)
	GetByName(name string) (models.Supplier, error)
	GetAll(limit int, offset int) ([]models.Supplier, error)
	Update(ID int, Input input.SupplierInput) (models.Supplier, error)
	Delete(ID int) (models.Supplier, error)
	ExportSuppliersToXLS() (*excelize.File, error)
	ImportSuppliersFromXLS(filePath string) ([]models.Supplier, error)
}

type supplierService struct {
	repository repository.SupplierRepository
}

func NewSupplierService(repository repository.SupplierRepository) *supplierService {
	return &supplierService{repository}
}

func (s *supplierService) CreateSupplier(input input.SupplierInput) (models.Supplier, error) {
	supplier := models.Supplier{
		Name:    input.Name,
		Address: input.Address,
		Email:   input.Email,
		Phone:   input.Phone,
	}

	rand.Seed(uint64(time.Now().UnixNano()))
	supplier.Code = rand.Intn(90000) + 10000

	if err := helper.ValidateEmail(supplier.Email); err != nil {
		return models.Supplier{}, err
	}

	if err := helper.ValidatePhoneNumber(supplier.Phone); err != nil {
		return models.Supplier{}, err
	}

	newSupplier, err := s.repository.Save(supplier)
	if err != nil {
		return newSupplier, err
	}

	return newSupplier, nil
}

func (s *supplierService) GetByID(ID int) (models.Supplier, error) {
	supplier, err := s.repository.FindByID(ID)
	if err != nil {
		return supplier, err
	}
	return supplier, nil
}

func (s *supplierService) GetByName(name string) (models.Supplier, error) {
	supplier, err := s.repository.FindByName(name)
	if err != nil {
		return supplier, err
	}
	return supplier, nil
}

func (s *supplierService) GetAll(limit int, offset int) ([]models.Supplier, error) {
	suppliers, err := s.repository.FindAll(limit, offset)
	if err != nil {
		return suppliers, err
	}
	return suppliers, nil
}

func (s *supplierService) Update(ID int, input input.SupplierInput) (models.Supplier, error) {
	supplier, err := s.repository.FindByID(ID)
	if err != nil {
		return supplier, err
	}

	supplier.Name = input.Name
	supplier.Address = input.Address
	supplier.Email = input.Email
	supplier.Phone = input.Phone

	updatedSupplier, err := s.repository.Update(ID, supplier)
	if err != nil {
		return updatedSupplier, err
	}
	return updatedSupplier, nil
}

func (s *supplierService) Delete(ID int) (models.Supplier, error) {
	supplier, err := s.repository.FindByID(ID)
	if err != nil {
		return supplier, err
	}
	deletedSupplier, err := s.repository.Delete(ID)
	if err != nil {
		return deletedSupplier, err
	}
	return deletedSupplier, nil
}

func (s *supplierService) ExportSuppliersToXLS() (*excelize.File, error) {
	// Fetch suppliers from the database
	suppliers, err := s.repository.FindAll(0, 0)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Number of suppliers retrieved: %d\n", len(suppliers))

	if len(suppliers) == 0 {
		return nil, errors.New("no suppliers found to export")
	}

	// Create a new Excel file
	f := excelize.NewFile()
	sheet := "Suppliers"
	index, err := f.NewSheet(sheet)
	if err != nil {
		return nil, err
	}

	// Write headers
	headers := []string{"ID", "Name", "Address", "Email", "Phone", "Code", "Created At", "Updated At"}
	for i, header := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheet, cell, header)
	}

	// Write supplier data
	for i, supplier := range suppliers {
		row := i + 2
		fmt.Printf("Writing supplier %d to row %d\n", supplier.ID, row) // Debugging output

		f.SetCellValue(sheet, "A"+strconv.Itoa(row), supplier.ID)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), supplier.Name)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), supplier.Address)
		f.SetCellValue(sheet, "D"+strconv.Itoa(row), supplier.Email)
		f.SetCellValue(sheet, "E"+strconv.Itoa(row), supplier.Phone)
		f.SetCellValue(sheet, "F"+strconv.Itoa(row), supplier.Code)
		f.SetCellValue(sheet, "G"+strconv.Itoa(row), supplier.CreatedAt.Format(time.RFC3339))
		f.SetCellValue(sheet, "H"+strconv.Itoa(row), supplier.UpdatedAt.Format(time.RFC3339))
	}

	// Set active sheet and delete the default sheet
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	return f, nil
}

func (s *supplierService) ImportSuppliersFromXLS(filePath string) ([]models.Supplier, error) {
	// Open the Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read data from the first sheet
	sheet := "Suppliers"
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}

	var importedSuppliers []models.Supplier

	// Skip the header row
	for i, row := range rows {
		if i == 0 {
			continue
		}

		// Parse each column into supplier fields
		code, _ := strconv.Atoi(row[5])

		supplier := models.Supplier{
			Name:    row[1],
			Address: row[2],
			Email:   row[3],
			Phone:   row[4],
			Code:    code,
		}

		// Insert supplier into the database
		savedSupplier, err := s.repository.Save(supplier)
		if err != nil {
			return nil, err
		}

		// Add the saved supplier to the list
		importedSuppliers = append(importedSuppliers, savedSupplier)
	}

	return importedSuppliers, nil
}
