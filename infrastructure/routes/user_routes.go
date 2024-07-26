package routes

import (
	"github.com/gin-gonic/gin"
	"pm/infrastructure/controllers/handlers"
)

type UserRoutes struct {
	handler *handlers.UserHandler
}

func NewUserRoutes(handler *handlers.UserHandler) *UserRoutes {

	return &UserRoutes{handler}
}

func (router *UserRoutes) RegisterRoutes(routerGroup *gin.RouterGroup) {
	users := routerGroup.Group("/users")
	{
		users.POST("/authenticate", router.handler.HandleAuthenticate)
		users.GET("", router.handler.HandleGetAllUsers)
		users.GET("/:id", router.handler.HandleGetUserByID)
		users.POST("", router.handler.HandleCreateUser)
		users.PUT("/:id", router.handler.HandleUpdateUserByID)
		users.DELETE("/:id", router.handler.HandleDeleteUserByID)

	}
}