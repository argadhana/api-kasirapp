package input

type TransactionProductInput struct {
	ProductID int `json:"product_id"`
	Qty       int `json:"quantity"`
}

type TransactionInput struct {
	Products []TransactionProductInput `json:"products"`
	Balance  float32                   `json:"balance"`
}
