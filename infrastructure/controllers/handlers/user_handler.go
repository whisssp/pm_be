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

// Authenticate 			godoc
// @Summary 			Authenticate user to get access resource
// @Description			Authenticate to receive a token string to use it for verifying permission
// Tag					User
// @Param				payload.LoginRequest body payload.LoginRequest{} true "payload.LoginRequest"
// @Success				200		{object} payload.AppResponse{}
// @Failure      		400  	{object} payload.AppError{}
// @Failure 			500 	{object} payload.AppError{}
// @Router				/users/authenticate [post]
func (h *UserHandler) HandleAuthenticate(c *gin.Context) {
	span := h.p.Logger.Start(c, "handlers/HandleAuthenticate")
	h.p.Logger.SetContextWithSpan(span)
	defer h.p.Logger.End()
	var loginRequest payload.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequest(err))
		return
	}

	response, err := h.userUsecase.Authenticate(c, &loginRequest)
	if err != nil {
		h.p.Logger.Error(err.Error(), map[string]interface{}{
			"message": err.Error(),
		})
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequest(err))
		return
	}

	h.p.Logger.Info("handlers/HandleAuthenticate", map[string]interface{}{"token": response.Token})
	c.JSON(http.StatusOK, payload.SuccessResponse(response, ""))
}

// Create User 			godoc
// @Summary 			Create a user
// @Description			Create a user to get info to authenticate
// Tag					User
// @Param				payload.UserRequest body payload.UserRequest{} true "payload.UserRequest"
// @Success				200		{object} payload.AppResponse{}
// @Failure      		400  	{object} payload.AppError{}
// @Failure 			500 	{object} payload.AppError{}
// @Router				/users/authenticate [post]
func (h *UserHandler) HandleCreateUser(c *gin.Context) {
	var userRequest payload.UserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequest(err))
		return
	}

	if err := h.userUsecase.CreateUser(c, &userRequest); err != nil {
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequest(err))
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