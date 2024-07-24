package entity

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name     string    `gorm:"type:varchar(255)"`
	Products []Product `gorm:"foreignKey:CategoryID"`
}