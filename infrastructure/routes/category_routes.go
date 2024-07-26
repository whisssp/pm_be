package routes

import (
	"github.com/gin-gonic/gin"
	"pm/infrastructure/controllers/handlers"
)

type CategoryRoutes struct {
	handler *handlers.CategoryHandler
}

func NewCategoryRoutes(handler *handlers.CategoryHandler) *CategoryRoutes {
	return &CategoryRoutes{handler}
}

func (router *CategoryRoutes) RegisterRoutes(routerGroup *gin.RouterGroup) {
	categories := routerGroup.Group("/categories")
	{
		categories.POST("", router.handler.HandleCreateCategory)
		//categories.GET("/search", router.handler.HandleGetAllCategories)
		categories.GET("", router.handler.HandleGetAllCategories)
		categories.GET("/:id", router.handler.HandleGetCategoryByID)
		categories.DELETE("/:id", router.handler.HandleDeleteCategoryByID)
		categories.PUT("/:id", router.handler.HandleUpdateCategoryByID)
	}
}