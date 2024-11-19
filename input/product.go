package input

type ProductInput struct {
	Name         string  `json:"name" validate:"required"`
	ProductType  string  `json:"product_type" validate:"required"`
	ImageURL     string  `form:"image_product" validate:"required"`
	BasePrice    float64 `json:"base_price" validate:"required"`
	SellingPrice float64 `json:"selling_price" validate:"required"`
	Stock        int     `json:"stock" validate:"required"`
	CodeProduct  string  `json:"code_product" validate:"required"`
	CategoryID   int     `json:"category_id" validate:"required"`
	MinimumStock int     `json:"minimum_stock" validate:"required"`
	Shelf        string  `json:"shelf" validate:"required"`
	Weight       int     `json:"weight" validate:"required"`
	Discount     int     `json:"discount" validate:"required"`
	Information  string  `json:"information" validate:"required"`
}
