package users

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"pm/domain/entity"
	"pm/domain/repository/users"
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

func NewUserRepository(c *gin.Context, p *base.Persistence, db *gorm.DB) users.UserRepository {
	return UserRepository{db, p, c}
}

func (u UserRepository) GetUserByEmail(parentSpan trace.Span, email string) (*entity.User, error) {
	span := u.p.Logger.Start(u.c, "GET_USER_BY_EMAIL: DATABASE", u.p.Logger.UseGivenSpan(parentSpan))
	defer span.End()
	u.p.Logger.Info("STARTING: GET USER BY EMAIL", map[string]interface{}{"email": email})

	db := u.db
	var user entity.User
	if err := db.Omit("Orders").Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			u.p.Logger.Error("GET_USER_BY_EMAIL: EMAIL DOESN'T EXISTS", map[string]interface{}{"error": entity.ErrEmailNotFound}, u.p.Logger.UseGivenSpan(span))
			return nil, payload.ErrEntityNotFound(entityName, entity.ErrEmailNotFound)
		}
		u.p.Logger.Error("GET_USER_BY_EMAIL: ERROR", map[string]interface{}{"error": err.Error()}, u.p.Logger.UseGivenSpan(span))
		return nil, payload.ErrDB(err)
	}
	u.p.Logger.Info("GET_USER_BY_EMAIL: SUCCESSFULLY", map[string]interface{}{"user": user}, u.p.Logger.UseGivenSpan(span))
	fmt.Println("SPAN-REPO", span)
	return &user, nil
}

func (u UserRepository) Create(parentSpan trace.Span, user *entity.User) error {
	span := u.p.Logger.Start(u.c, "CREATE_USER: DATABASE", u.p.Logger.UseGivenSpan(parentSpan))
	defer span.End()
	u.p.Logger.Info("STARTING: CREATE USER", map[string]interface{}{"data": user}, u.p.Logger.UseGivenSpan(span))

	db := u.db
	if err := db.Debug().Model(&entity.User{}).Create(user).Error; err != nil {
		u.p.Logger.Info("CREATE_USER: ERROR", map[string]interface{}{"error": err.(*pgconn.PgError).Detail}, u.p.Logger.UseGivenSpan(span))
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return payload.ErrInvalidRequest(fmt.Errorf("%s", err.(*pgconn.PgError).Detail))
		}
		return payload.ErrDB(err)
	}
	u.p.Logger.Info("CREATE_USER: SUCCESSFULLY", map[string]interface{}{"user": user}, u.p.Logger.UseGivenSpan(span))
	return nil
}

func (u UserRepository) Update(parentSpan trace.Span, user *entity.User) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserRepository) GetAllUsers(parentSpan trace.Span) ([]entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserRepository) GetUserByID(parentSpan trace.Span, id int64) (*entity.User, error) {
	span := u.p.Logger.Start(u.c, "GET_USER_BY_ID: DATABASE", u.p.Logger.UseGivenSpan(parentSpan))
	defer span.End()
	u.p.Logger.Info("STARTING: GET USER BY ID", map[string]interface{}{"id": id}, u.p.Logger.UseGivenSpan(span))

	db := u.db
	var user entity.User
	if err := db.Where("id = ?", id).Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			u.p.Logger.Info("GET_USER_BY_ID: USER ID DOESN'T EXISTS", map[string]interface{}{"error": err.Error()}, u.p.Logger.UseGivenSpan(span))
			return nil, payload.ErrEntityNotFound("users", err)
		}
		u.p.Logger.Info("GET_USER_BY_ID: ERROR", map[string]interface{}{"error": err.Error()}, u.p.Logger.UseGivenSpan(span))
		return nil, payload.ErrDB(err)
	}
	u.p.Logger.Info("GET_USER_BY_ID_SUCCESSFULLY", map[string]interface{}{"user": user}, u.p.Logger.UseGivenSpan(span))
	return &user, nil
}

func (u UserRepository) GetUserByRole(parentSpan trace.Span, role entity.UserRole) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserRepository) Delete(parentSpan trace.Span, user *entity.User) error {
	//TODO implement me
	panic("implement me")
}