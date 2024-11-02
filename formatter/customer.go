package formatter

import "api-kasirapp/models"

type CustomerFormatter struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func FormatCustomer(customer models.Customer) CustomerFormatter {
	formatter := CustomerFormatter{
		ID:        customer.ID,
		Name:      customer.Name,
		Address:   customer.Address,
		Phone:     customer.Phone,
		Email:     customer.Email,
		CreatedAt: customer.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: customer.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return formatter
}

func FormatCustomers(customers []models.Customer) []CustomerFormatter {
	var customersFormatter []CustomerFormatter
	for _, customer := range customers {
		formatter := FormatCustomer(customer)
		customersFormatter = append(customersFormatter, formatter)
	}
	return customersFormatter
}
