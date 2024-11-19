package models

import "time"

type Shift struct {
	ID           int       `gorm:"primaryKey;autoIncrement"`
	UserID       User      `gorm:"foreignKey:ShiftName;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;not null"`
	StartBalance float64   `gorm:"not null"`
	StartTime    time.Time `gorm:"not null"`
	EndTime      *time.Time
	Status       string `gorm:"default:berjalan"`
	TotalSales   float64
	Expenses     float64
	CreatedAt    string
	UpdatedAt    string
}
