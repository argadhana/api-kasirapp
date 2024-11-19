package formatter

import "api-kasirapp/models"

type TransactionDetailFormatter struct {
	ProductID int `json:"product_id"`
	Qty       int `json:"qty"`
}

type TransactionFormatter struct {
	ID         int                          `json:"id"`
	Details    []TransactionDetailFormatter `json:"details"`
	Amount     float64                      `json:"amount"`
	CashReturn float64                      `json:"cash_return"`
	CreatedAt  string                       `json:"created_at"`
	UpdatedAt  string                       `json:"updated_at"`
}

func FormatTransaction(transaction models.Transaction, cashReturn float64) TransactionFormatter {
	var details []TransactionDetailFormatter
	for _, detail := range transaction.Details {
		details = append(details, TransactionDetailFormatter{
			ProductID: detail.ProductID,
			Qty:       detail.Qty,
		})
	}

	formatter := TransactionFormatter{
		ID:         transaction.ID,
		Details:    details,
		Amount:     transaction.Amount,
		CashReturn: cashReturn,
		CreatedAt:  transaction.CreatedAt.String(),
		UpdatedAt:  transaction.UpdatedAt.String(),
	}
	return formatter
}
