package formatter

import "api-kasirapp/models"

type TransactionFormatter struct {
	ID        int     `json:"id"`
	ProductID int     `json:"product_id"`
	Qty       int     `json:"qty"`
	Amount    float32 `json:"amount"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type TransactionFormatterWithProducts struct {
	ID        int              `json:"id"`
	ProductID int              `json:"product_id"`
	Qty       int              `json:"qty"`
	Amount    float32          `json:"amount"`
	Product   ProductFormatter `json:"product"`
	CreatedAt string           `json:"created_at"`
	UpdatedAt string           `json:"updated_at"`
}

func FormatTransaction(transaction *models.Transaction) TransactionFormatter {
	formatter := TransactionFormatter{
		ID:        transaction.ID,
		ProductID: transaction.ProductID,
		Qty:       transaction.Qty,
		Amount:    transaction.Amount,
		CreatedAt: transaction.CreatedAt.String(),
		UpdatedAt: transaction.UpdatedAt.String(),
	}
	return formatter
}

func FormatTransactionsWithProducts(transactions []models.Transaction) []TransactionFormatterWithProducts {
	formatters := []TransactionFormatterWithProducts{}

	for _, transaction := range transactions {
		formatter := TransactionFormatterWithProducts{
			ID:        transaction.ID,
			ProductID: transaction.ProductID,
			Qty:       transaction.Qty,
			Amount:    transaction.Amount,
			Product:   FormatProduct(transaction.Product),
			CreatedAt: transaction.CreatedAt.String(),
			UpdatedAt: transaction.UpdatedAt.String(),
		}

		formatters = append(formatters, formatter)
	}

	return formatters
}
