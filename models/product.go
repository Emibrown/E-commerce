package models

import "time"

type Product struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	Price       float64 `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
