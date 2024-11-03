package models

import "time"

type Transaction struct {
	ID        int `gorm:"primaryKey;autoIncrement"`
	ProductID int `gorm:"index;column:id_product"`
	Qty       int
	Amount    float32
	CreatedAt time.Time
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Product   Product   `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
