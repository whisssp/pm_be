package entity

import "gorm.io/gorm"

// PaymentID int64
type Order struct {
	gorm.Model
	UserID     uint
	Status     string      `gorm:"type:varchar(50)"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint
	ProductID uint
	Quantity  int
	Price     float64 `gorm:"type:double precision"`
}