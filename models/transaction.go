package models

import "time"

type Transaction struct {
	ID        int                 `gorm:"primaryKey;autoIncrement" json:"id"`
	Qty       int                 `gorm:"not null" json:"quantity"`                                            // Total number of items in the transaction
	Amount    float64             `gorm:"not null" json:"amount"`                                              // Total amount for the transaction
	Details   []TransactionDetail `gorm:"foreignKey:TransactionID;constraint:OnDelete:CASCADE" json:"details"` // Associated transaction details
	CreatedAt time.Time           `gorm:"autoCreateTime" json:"created_at"`                                    // Automatically set on creation
	UpdatedAt time.Time           `gorm:"autoUpdateTime" json:"updated_at"`                                    // Automatically updated on modification
}

type TransactionDetail struct {
	ID            int       `gorm:"primaryKey;autoIncrement" json:"id"`
	TransactionID int       `gorm:"not null;index" json:"transaction_id"`  // Foreign key to transactions
	ProductID     int       `gorm:"not null;index" json:"product_id"`      // Foreign key to products
	Qty           int       `gorm:"not null" json:"quantity"`              // Quantity of the product in the transaction
	Product       Product   `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"product"` // Associated product
}