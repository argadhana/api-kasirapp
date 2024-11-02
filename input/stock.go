package input

type StockInput struct {
	Date        string  `json:"date"`
	BuyingPrice float64 `json:"buying_price"`
	Amount      int     `json:"amount"`
	Information string  `json:"information"`
}
