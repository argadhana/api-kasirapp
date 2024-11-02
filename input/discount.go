package input

type DiscountInput struct {
	Name       string  `json:"name" form:"name" validate:"required"`
	Percentage float64 `json:"percentage" form:"percentage" validate:"required"`
}
