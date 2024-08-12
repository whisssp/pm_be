package application

import (
	"errors"
	"github.com/gin-gonic/gin"
	"pm/domain/entity"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/implementations/users"
	"pm/infrastructure/mapper"
	"pm/infrastructure/persistences/base"
	"pm/infrastructure/persistences/base/logger"
	"pm/utils"
)

type UserUsecase interface {
	CreateUser(*gin.Context, *payload.UserRequest) error
	GetUserByID(id int64) (*entity.User, error)
	GetAllUsers() ([]payload.ListUserResponses, error)
	UpdateUserByID(request *payload.UserRequest) (*entity.User, error)
	DeleteUserByID(id int64) error
	Authenticate(*gin.Context, *payload.LoginRequest) (string, error)
}

type userUsecase struct {
	p *base.Persistence
}

func NewUserUsecase(p *base.Persistence) UserUsecase {
	return userUsecase{p}
}

func (u userUsecase) Authenticate(c *gin.Context, request *payload.LoginRequest) (string, error) {
	newlogger := logger.NewLogger()
	ctx, span := newlogger.Start(c, "AUTHENTICATE_USECASES")
	defer newlogger.End()
	newlogger.Info("STARTING: AUTHENTICATE", map[string]interface{}{"data": request})

	if err := utils.ValidateReqPayload(request); err != nil {
		newlogger.Error("AUTHENTICATE: INVALID REQUEST", map[string]interface{}{"message": err.Error()})
		return "", payload.ErrInvalidRequest(err)
	}
	userRepo := users.NewUserRepository(ctx, u.p, u.p.GormDB)

	user, err := userRepo.GetUserByEmail(request.Email)
	if err != nil {
		newlogger.Error("AUTHENTICATE: EMAIL DOESN'T EXISTS", map[string]interface{}{"error": err.Error()})
		return "", err
	}

	if !utils.ComparePasswords([]byte(user.Password), []byte(request.Password)) {
		//cspan := newlogger.Start(c, "AUTHENTICATE: PASSWORD FAILED")
		//defer cspan.End()
		newlogger.Error("AUTHENTICATE: WRONG PASSWORD", map[string]interface{}{"error": "password is incorrect"})
		return "", payload.ErrWrongPassword(errors.New("incorrect email or password"))
	}

	token, err := utils.JwtGenerateJwtToken(ctx, u.p, user, span)
	if err != nil {
		newlogger.Error("AUTHENTICATE: GENERATE TOKEN FAILED", map[string]interface{}{"error": err.Error()})
		return "", err
	}

	newlogger.Info("AUTHENTICATE: SUCCESSFULLY", map[string]interface{}{"authenticate_response": token})
	return token, nil
}

func (u userUsecase) CreateUser(c *gin.Context, request *payload.UserRequest) error {
	_, _ = u.p.Logger.Start(c, "CREATE_USER_USECASES")
	defer u.p.Logger.End()
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
		return err
	}
	u.p.Logger.Info("CREATE_USER_SUCCESSFULLY", map[string]interface{}{"data": user})
	return nil
}

func (u userUsecase) GetUserByID(id int64) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u userUsecase) GetAllUsers() ([]payload.ListUserResponses, error) {
	//TODO implement me
	panic("implement me")
}

func (u userUsecase) UpdateUserByID(request *payload.UserRequest) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u userUsecase) DeleteUserByID(id int64) error {
	//TODO implement me
	panic("implement me")
}
