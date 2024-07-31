package application

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/implementations/users"
	"pm/infrastructure/mapper"
	"pm/infrastructure/persistences/base"
	"pm/utils"
)

type UserUsecase interface {
	CreateUser(*gin.Context, *payload.UserRequest) error
	GetUserByID(id int64) (*payload.UserResponse, error)
	GetAllUsers() ([]payload.ListUserResponses, error)
	UpdateUserByID(request *payload.UserRequest) (*payload.UserResponse, error)
	DeleteUserByID(id int64) error
	Authenticate(*gin.Context, *payload.LoginRequest) (*payload.AuthResponse, error)
}

type userUsecase struct {
	p *base.Persistence
}

func NewUserUsecase(p *base.Persistence) UserUsecase {
	return userUsecase{p}
}

func (u userUsecase) Authenticate(c *gin.Context, request *payload.LoginRequest) (*payload.AuthResponse, error) {
	//channels := []string{"Honeycomb"}

	if err := utils.ValidateReqPayload(request); err != nil {
		return nil, payload.ErrInvalidRequest(err)
	}

	userRepo := users.NewUserRepository(u.p.GormDB)
	user, err := userRepo.GetUserByEmail(request.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, payload.ErrEntityNotFound(user.TableName(), err)
		}
		return nil, payload.ErrDB(err)
	}

	if !utils.ComparePasswords([]byte(user.Password), []byte(request.Password)) {
		return nil, payload.ErrWrongPassword(errors.New("incorrect email or password"))
	}

	token, err := utils.JwtGenerateJwtToken(user)
	if err != nil {
		return nil, payload.ErrGenerateToken(err)
	}

	authResponse := payload.AuthResponse{Token: token}
	return &authResponse, nil
}

func (u userUsecase) CreateUser(c *gin.Context, request *payload.UserRequest) error {
	if err := utils.ValidateReqPayload(request); err != nil {
		return payload.ErrInvalidRequest(err)
	}
	userRepo := users.NewUserRepository(u.p.GormDB)
	user := mapper.UserRequestToUser(request)
	hashed, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	if !utils.ComparePasswords([]byte(hashed), []byte(user.Password)) {
		return errors.New("invalid hash")
	}
	user.Password = hashed
	if err := userRepo.Create(user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return payload.ErrExisted(err)
		}

		return payload.ErrDB(err)
	}
	return nil
}

func (u userUsecase) GetUserByID(id int64) (*payload.UserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u userUsecase) GetAllUsers() ([]payload.ListUserResponses, error) {
	//TODO implement me
	panic("implement me")
}

func (u userUsecase) UpdateUserByID(request *payload.UserRequest) (*payload.UserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u userUsecase) DeleteUserByID(id int64) error {
	//TODO implement me
	panic("implement me")
}