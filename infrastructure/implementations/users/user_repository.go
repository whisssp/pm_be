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
	span := u.p.Logger.Start(u.c, "GET_USER_BY_EMAIL: DATABASE", u.p.Logger.SetContextWithSpanFunc())
	defer span.End()
	u.p.Logger.Info("STARTING: GET USER BY EMAIL", map[string]interface{}{"email": email}, u.p.Logger.SetContextWithSpanFunc())

	db := u.db
	var user entity.User
	if err := db.Omit("Orders").Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			u.p.Logger.Error("GET_USER_BY_EMAIL: EMAIL DOESN'T EXISTS", map[string]interface{}{"error": entity.ErrEmailNotFound}, u.p.Logger.SetContextWithSpanFunc())
			return nil, payload.ErrEntityNotFound(entityName, entity.ErrEmailNotFound)
		}
		u.p.Logger.Error("GET_USER_BY_EMAIL: ERROR", map[string]interface{}{"error": err.Error()}, u.p.Logger.SetContextWithSpanFunc())
		return nil, payload.ErrDB(err)
	}
	u.p.Logger.Info("GET_USER_BY_EMAIL: SUCCESSFULLY", map[string]interface{}{"user": user}, u.p.Logger.SetContextWithSpanFunc())
	return &user, nil
}

func (u UserRepository) Create(user *entity.User) error {
	span := u.p.Logger.Start(u.c, "CREATE_USER: DATABASE")
	defer span.End()
	u.p.Logger.Info("STARTING: CREATE USER", map[string]interface{}{"data": user})

	db := u.db
	if err := db.Debug().Model(&entity.User{}).Create(user).Error; err != nil {
		u.p.Logger.Info("CREATE_USER: ERROR", map[string]interface{}{"error": err.Error()})
		if errors.Is(err, gorm.ErrInvalidData) {
			return payload.ErrInvalidRequest(err)
		}
		return payload.ErrDB(err)
	}
	u.p.Logger.Info("CREATE_USER: SUCCESSFULLY", map[string]interface{}{"user": user})
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
	span := u.p.Logger.Start(u.c, "GET_USER_BY_ID: DATABASE")
	defer span.End()
	u.p.Logger.Info("STARTING: GET USER BY ID", map[string]interface{}{"id": id})

	db := u.db
	var user entity.User
	if err := db.Where("id = ?", id).Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			u.p.Logger.Info("GET_USER_BY_ID: USER ID DOESN'T EXISTS", map[string]interface{}{"error": err.Error()})
			return nil, payload.ErrEntityNotFound("users", err)
		}
		u.p.Logger.Info("GET_USER_BY_ID: ERROR", map[string]interface{}{"error": err.Error()})
		return nil, payload.ErrDB(err)
	}
	u.p.Logger.Info("GET_USER_BY_ID_SUCCESSFULLY", map[string]interface{}{"user": user})
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