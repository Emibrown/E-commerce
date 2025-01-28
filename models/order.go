package models

import (
	"time"
)

type OrderStatus string

const (
	Pending   OrderStatus = "Pending"
	Shipped   OrderStatus = "Shipped"
	Completed OrderStatus = "Completed"
	Cancelled OrderStatus = "Cancelled"
)

type Order struct {
	ID        uint        `gorm:"primaryKey"`
	UserID    uint        `gorm:"not null"`
	Products  []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	Status    OrderStatus `gorm:"type:varchar(20); default:'Pending'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrderItem struct {
	ID        uint    `gorm:"primaryKey"`
	OrderID   uint    `gorm:"not null"`
	ProductID uint    `gorm:"not null"`
	Product   Product `gorm:"foreignKey:ProductID"`
	Quantity  int     `gorm:"not null; default:1"`
	Price     float64 `gorm:"not null"` // capture price at the time of ordering
}
