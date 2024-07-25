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
	category := routerGroup.Group("/categories")
	{
		category.POST("", router.handler.HandleCreateCategory)
		//categories.GET("/search", router.handler.HandleGetAllCategories)
		category.GET("", router.handler.HandleGetAllCategories)
		category.GET("/:id", router.handler.HandleGetCategoryByID)
		category.DELETE("/:id", router.handler.HandleDeleteCategoryByID)
		category.PUT("/:id", router.handler.HandleUpdateCategoryByID)
	}
}