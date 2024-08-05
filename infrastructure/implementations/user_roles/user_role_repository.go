package user_roles

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"pm/domain/entity"
	"pm/domain/repository/user_roles"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/persistences/base"
)

type UserRoleRepository struct {
	db *gorm.DB
	p  *base.Persistence
	c  *gin.Context
}

func (u UserRoleRepository) GetUserRoleByID(i int64) (*entity.UserRole, error) {
	//u.p.Logger.SetContextWithSpan(span)
	span := u.p.Logger.Start(u.c, "GET_USER_ROLE_BY_ID: DATABASE")
	defer span.End()
	u.p.Logger.Info("STARTING: GET USER ROLE BY ID", map[string]interface{}{"id": i}, u.p.Logger.SetContextWithSpanFunc())
	var userRole entity.UserRole
	if err := u.db.Debug().Model(&entity.UserRole{}).Where("id = ?", i).First(&userRole).Error; err != nil {
		u.p.Logger.Error("GET_USER_BY_ID: ERROR", map[string]interface{}{"error": err.Error()}, u.p.Logger.SetContextWithSpanFunc())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, payload.ErrEntityNotFound("user_roles", err)
		}
		return nil, payload.ErrDB(err)
	}
	u.p.Logger.Info("GET_USER_BY_ID: SUCCESSFULLY", map[string]interface{}{"user_role": userRole}, u.p.Logger.SetContextWithSpanFunc())
	return &userRole, nil
}

func (u UserRoleRepository) GetUserRoleByName(s string) (*entity.UserRole, error) {
	span := u.p.Logger.Start(u.c, "GET_USER_ROLE_BY_NAME: DATABASE")
	defer span.End()
	u.p.Logger.Info("STARTING: GET USER ROLE BY NAME", map[string]interface{}{"role_name": s})

	var userRole entity.UserRole
	if err := u.db.Debug().Model(&entity.UserRole{}).Where("name = ?", s).First(&userRole).Error; err != nil {
		u.p.Logger.Error("GET_USER_ROLE_BY_NAME: ERROR", map[string]interface{}{"error": err.Error()})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, payload.ErrEntityNotFound("user_roles", err)
		}
		return nil, payload.ErrDB(err)
	}
	u.p.Logger.Info("GET_USER_ROLE_BY_NAME: SUCCESSFULLY", map[string]interface{}{"user_role": userRole})
	return &userRole, nil
}

func NewUserRoleRepository(db *gorm.DB, p *base.Persistence, c *gin.Context) user_roles.UserRoleRepository {
	return UserRoleRepository{
		db: db,
		p:  p,
		c:  c,
	}
}