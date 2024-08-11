package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pm/application"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/persistences/base"
)

type UserHandler struct {
	p           *base.Persistence
	userUsecase application.UserUsecase
}

func NewUserHandler(p *base.Persistence) *UserHandler {
	userUsecase := application.NewUserUsecase(p)
	return &UserHandler{p, userUsecase}
}

// HandleAuthenticate 	Authenticate 			godoc
// @Summary 			Authenticate user to get access resource
// @Description			Authenticate to receive a token string to use it for verifying permission
// @Tags				User
// @Accept				json
// @Produce				json
// @Param				LoginRequest body payload.LoginRequest true "send the login request data to authenticate"
// @Success				200		{object} payload.AppResponse
// @Failure      		400  	{object} payload.AppError
// @Failure 			500 	{object} payload.AppError
// @Router				/users/authenticate [post]
func (h *UserHandler) HandleAuthenticate(c *gin.Context) {
	span := h.p.Logger.Start(c, "handlers/HandleAuthenticate", h.p.Logger.SetContextWithSpanFunc())
	defer span.End()

	var loginRequest payload.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.Error(payload.ErrInvalidRequest(err))
		h.p.Logger.Error("AUTHENTICATE_FAILED", map[string]interface{}{"message": err.Error()})
		return
	}

	token, err := h.userUsecase.Authenticate(c, &loginRequest)
	if err != nil {
		c.Error(err)
		h.p.Logger.Error("AUTHENTICATE_FAILED", map[string]interface{}{
			"message": err.Error(),
		})
		return
	}

	authResponse := payload.AuthResponse{Token: token}
	c.JSON(http.StatusOK, payload.SuccessResponse(authResponse, ""))
}

// HandleCreateUser Create User 			godoc
// @Summary 			Create a user
// @Description			Create a user to get info to authenticate
// @Tags				User
// @Accept				json
// @Produce				json
// @Param				UserRequest body payload.UserRequest true "create new user"
// @Success				200		{object} payload.AppResponse
// @Failure      		400  	{object} payload.AppError
// @Failure 			500 	{object} payload.AppError
// @Router				/users [post]
func (h *UserHandler) HandleCreateUser(c *gin.Context) {
	span := h.p.Logger.Start(c, "handlers/HandleCreateUser", h.p.Logger.SetContextWithSpanFunc())
	defer span.End()

	var userRequest payload.UserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.Error(err)
		h.p.Logger.Error("CREATE_USER_FAILED", map[string]interface{}{
			"message": err.Error(),
		})
		return
	}

	if err := h.userUsecase.CreateUser(c, &userRequest); err != nil {
		c.Error(err)
		h.p.Logger.Error("CREATE_USER_FAILED", map[string]interface{}{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, payload.SuccessResponse(nil, ""))
}

func (h *UserHandler) HandleGetUserByID(c *gin.Context) {

}

func (h *UserHandler) HandleUpdateUserByID(c *gin.Context) {

}

func (h *UserHandler) HandleGetAllUsers(c *gin.Context) {

}

func (h *UserHandler) HandleDeleteUserByID(c *gin.Context) {

}