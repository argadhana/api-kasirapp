package input

type CategoryInput struct {
	Name string `json:"name" form:"name" validate:"required"`
}
