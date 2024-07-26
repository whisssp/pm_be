package entity

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(150)"`
	Phone    string `gorm:"type:varchar(11)"`
	Password string `gorm:"type:varchar(255)"`
	Role     string `gorm:"type:varchar(50);default:'USER'"`
	Email    string `gorm:"type:varchar(200);unique"`
}

func (u *User) TableName() string {
	return "users"
}