package routes

import (
	"github.com/gin-gonic/gin"
	"pm/infrastructure/controllers/handlers"
	"pm/infrastructure/controllers/middleware"
)

type ProductRoutes struct {
	handler *handlers.ProductHandler
}

func NewProductRoutes(handler *handlers.ProductHandler) *ProductRoutes {
	return &ProductRoutes{handler}
}

func (router *ProductRoutes) RegisterRoutes(routerGroup *gin.RouterGroup) {
	products := routerGroup.Group("/products").Use(middleware.AuthMiddleware())
	{
		products.POST("", router.handler.HandleCreateProduct)
		products.GET("/search", router.handler.HandleGetAllProducts)
		products.GET("", router.handler.HandleGetAllProducts)
		products.GET("/:id", router.handler.HandleGetProductByID)
		products.DELETE("/:id", router.handler.HandleDeleteProductByID)
		products.PUT("/:id", router.handler.HandleUpdateProductByID)
	}
}