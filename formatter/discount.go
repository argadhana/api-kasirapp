package formatter

import (
	"api-kasirapp/models"
)

type DiscountFormatter struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Percentage float64 `json:"percentage"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

func FormatDiscount(discount models.Discount) DiscountFormatter {
	formatter := DiscountFormatter{
		ID:         discount.ID,
		Name:       discount.Name,
		Percentage: discount.Percentage,
		CreatedAt:  discount.CreatedAt.String(),
		UpdatedAt:  discount.UpdatedAt.String(),
	}
	return formatter
}

func FormatDiscounts(discounts []models.Discount) []DiscountFormatter {
	var formatter []DiscountFormatter
	for _, discount := range discounts {
		formatter = append(formatter, FormatDiscount(discount))
	}
	return formatter
}
