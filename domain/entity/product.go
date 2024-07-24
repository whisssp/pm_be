package entity

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string `gorm:"type:varchar(255)"`
	Description string
	Price       float64 `gorm:"type:double precision"`
	CategoryID  int64
	Stock       int64
	Image       string `gorm:"type:text"`
}

func GetID(p Product) int64 {
	return int64(uint64(p.ID))
}