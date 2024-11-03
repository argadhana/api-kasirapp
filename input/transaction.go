package input

type TransactionInput struct {
	ProductID int     `json:"product_id"`
	Qty       int     `json:"qty"`
	Amount    float32 `json:"amount"`
}
