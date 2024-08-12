package users

import (
	"pm/domain/entity"
)

type UserRepository interface {
	Create(*entity.User) error
	Update(*entity.User) (*entity.User, error)
	GetAllUsers() ([]entity.User, error)
	GetUserByID(int64) (*entity.User, error)
	GetUserByRole(entity.UserRole) (*entity.User, error)
	Delete(*entity.User) error
	GetUserByEmail(string) (*entity.User, error)
}
