package users

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"pm/domain/entity"
	"pm/domain/repository"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/persistences/base"
)

const (
	entityName string = "users"
)

type UserRepository struct {
	db *gorm.DB
	p  *base.Persistence
	c  *gin.Context
}

func NewUserRepository(c *gin.Context, p *base.Persistence, db *gorm.DB) repository.UserRepository {
	return UserRepository{db, p, c}
}

func (u UserRepository) GetUserByEmail(email string) (*entity.User, error) {
	span := u.p.Logger.Start(u.c, "GET_USER_BY_EMAIL_DATABASE")
	defer span.End()
	u.p.Logger.Info("GET_USER_BY_EMAIL", map[string]interface{}{"data": email})

	db := u.db
	var user entity.User
	if err := db.Omit("Orders").Where("email = ?", email).Find(&user).Error; err != nil {
		u.p.Logger.Error("GET_USER_BY_EMAIL_FAILED", map[string]interface{}{"message": err.Error()})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, payload.ErrEntityNotFound(entityName, err)
		}
		return nil, payload.ErrDB(err)
	}
	u.p.Logger.Info("GET_USER_BY_EMAIL_SUCCESSFULLY", map[string]interface{}{"data": user})
	return &user, nil
}

func (u UserRepository) Create(user *entity.User) error {
	span := u.p.Logger.Start(u.c, "CREATE_USER_DATABASE")
	defer span.End()
	u.p.Logger.Info("CREATE_USER", map[string]interface{}{"data": user})

	db := u.db
	if err := db.Debug().Model(&entity.User{}).Create(user).Error; err != nil {
		u.p.Logger.Info("CREATE_USER_FAILED", map[string]interface{}{"message": err.Error()})
		return payload.ErrDB(err)
	}
	u.p.Logger.Info("CREATE_USER_SUCCESSFULLY", map[string]interface{}{"data": user})
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
	span := u.p.Logger.Start(u.c, "GET_USER_BY_ID_DATABASE")
	defer span.End()
	u.p.Logger.Info("GET_USER_BY_ID", map[string]interface{}{"data": id})
	db := u.db
	var user entity.User
	if err := db.Where("id = ?", id).Find(&user).Error; err != nil {
		u.p.Logger.Info("GET_USER_BY_ID_FAILED", map[string]interface{}{"message": err.Error()})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, payload.ErrDB(err)
	}
	u.p.Logger.Info("GET_USER_BY_ID_SUCCESSFULLY", map[string]interface{}{"data": user})
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