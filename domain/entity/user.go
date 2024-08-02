package entity

import (
	"errors"
	"gorm.io/gorm"
)

var (
	ErrEmailNotFound = errors.New("email not found")
	ErrUserNotFound  = errors.New("user not found")
)

type User struct {
	gorm.Model
	Name     string  `gorm:"type:varchar(150)"`
	Phone    string  `gorm:"type:varchar(11)"`
	Password string  `gorm:"type:varchar(255)"`
	RoleID   int64   `gorm:"type:bigint"`
	Email    string  `gorm:"type:varchar(200);unique"`
	Orders   []Order `gorm:"foreignKey:UserID"`
}

func (u *User) TableName() string {
	return "users"
}