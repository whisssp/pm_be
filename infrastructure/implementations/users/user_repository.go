package users

import (
	"errors"
	"gorm.io/gorm"
	"pm/domain/entity"
	"pm/domain/repository"
	"pm/infrastructure/controllers/payload"
)

const (
	entityName string = "users"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return UserRepository{db}
}

func (u UserRepository) GetUserByEmail(email string) (*entity.User, error) {
	db := u.db
	var user entity.User
	if err := db.Omit("Orders").Where("email = ?", email).Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, payload.ErrEntityNotFound(entityName, err)
		}
		return nil, payload.ErrDB(err)
	}
	return &user, nil
}

func (u UserRepository) Create(user *entity.User) error {
	db := u.db
	if err := db.Create(user).Error; err != nil {
		return payload.ErrDB(err)
	}
	return nil
}

func (u UserRepository) Update(user *entity.User) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserRepository) GetAllUsers() ([]entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserRepository) GetUserByID(id int64) (*entity.User, error) {
	db := u.db
	var user entity.User
	if err := db.Where("id = ?", id).Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, payload.ErrDB(err)
	}
	return &user, nil
}

func (u UserRepository) GetUserByRole(role entity.UserRole) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserRepository) Delete(user *entity.User) error {
	//TODO implement me
	panic("implement me")
}