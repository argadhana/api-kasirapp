package input

type SupplierInput struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
	Phone   string `json:"phone" validate:"required,regexp=^08|628[0-9]{9,11}$"`
}
