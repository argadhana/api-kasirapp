package input

type CustomerInput struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	Phone   string `json:"phone" validate:"required, len=13, regexp=^08|628[0-9]{9,}$"`
	Email   string `json:"email" validate:"optional, email"`
}
