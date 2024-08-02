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
	span := u.p.Logger.Start(c, "AUTHENTICATE_USECASES", u.p.Logger.SetContextWithSpanFunc())
	defer span.End()
	u.p.Logger.Info("STARTING: AUTHENTICATE", map[string]interface{}{"data": request})

	if err := utils.ValidateReqPayload(request); err != nil {
		u.p.Logger.Error("AUTHENTICATE: INVALID REQUEST", map[string]interface{}{"message": err.Error()}, u.p.Logger.SetContextWithSpanFunc())
		return nil, payload.ErrInvalidRequest(err)
	}

	userRepo := users.NewUserRepository(c, u.p, u.p.GormDB)
	user, err := userRepo.GetUserByEmail(request.Email)
	if err != nil {
		u.p.Logger.Error("AUTHENTICATE: EMAIL DOESN'T EXISTS", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	if !utils.ComparePasswords([]byte(user.Password), []byte(request.Password)) {
		u.p.Logger.Error("AUTHENTICATE: WRONG PASSWORD", map[string]interface{}{"error": "password is incorrect"}, u.p.Logger.UseGivenSpan(span))
		return nil, payload.ErrWrongPassword(errors.New("incorrect email or password"))
	}

	token, err := utils.JwtGenerateJwtToken(c, u.p, user, span)
	if err != nil {
		u.p.Logger.Error("AUTHENTICATE: GENERATE TOKEN FAILED", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	authResponse := payload.AuthResponse{Token: token}
	u.p.Logger.Info("AUTHENTICATE: SUCCESSFULLY", map[string]interface{}{"authenticate_response": authResponse})
	return &authResponse, nil
}

func (u userUsecase) CreateUser(c *gin.Context, request *payload.UserRequest) error {
	span := u.p.Logger.Start(c, "CREATE_USER_USECASES", u.p.Logger.SetContextWithSpanFunc())
	defer span.End()
	u.p.Logger.Info("CREATE_USER", map[string]interface{}{"data": request})

	if err := utils.ValidateReqPayload(request); err != nil {
		u.p.Logger.Error("CREATE_USER_FAILED", map[string]interface{}{"message": err.Error()})
		return payload.ErrInvalidRequest(err)
	}
	userRepo := users.NewUserRepository(c, u.p, u.p.GormDB)
	user := mapper.UserRequestToUser(request)
	hashed, err := utils.HashPassword(user.Password)
	if err != nil {
		u.p.Logger.Error("CREATE_USER_FAILED", map[string]interface{}{"message": err.Error()})
		return err
	}
	if !utils.ComparePasswords([]byte(hashed), []byte(user.Password)) {
		err = errors.New("invalid hash")
		u.p.Logger.Error("CREATE_USER_FAILED", map[string]interface{}{"message": err.Error()})
		return err
	}

	u.p.Logger.Info("CREATE_USER_INFO", map[string]interface{}{"data": user})
	user.Password = hashed
	if err := userRepo.Create(user); err != nil {
		u.p.Logger.Error("CREATE_USER_FAILED", map[string]interface{}{"message": err.Error()})
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