package user_roles

import (
	"pm/domain/entity"
)

type UserRoleRepository interface {
	GetUserRoleByID(int64) (*entity.UserRole, error)
	GetUserRoleByName(string) (*entity.UserRole, error)
}