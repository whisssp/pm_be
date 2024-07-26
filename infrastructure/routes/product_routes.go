package routes

import (
	"github.com/gin-gonic/gin"
	"pm/domain/entity"
	"pm/infrastructure/controllers/handlers"
	"pm/infrastructure/controllers/middleware"
	"pm/infrastructure/persistences/base"
)

type ProductRoutes struct {
	handler *handlers.ProductHandler
	p       *base.Persistence
}

func NewProductRoutes(p *base.Persistence, handler *handlers.ProductHandler) *ProductRoutes {
	return &ProductRoutes{handler, p}
}

func (router *ProductRoutes) RegisterRoutes(routerGroup *gin.RouterGroup) {
	products := routerGroup.Group("/products")
	{
		products.POST("", middleware.AuthMiddleware(router.p, entity.AllRoles[0]), router.handler.HandleCreateProduct)
		products.GET("/search", middleware.AuthMiddleware(router.p, entity.AllRoles[0]), router.handler.HandleGetAllProducts)
		products.GET("", router.handler.HandleGetAllProducts)
		products.GET("/:id", middleware.AuthMiddleware(router.p, entity.AllRoles[1]), router.handler.HandleGetProductByID)
		products.DELETE("/:id", router.handler.HandleDeleteProductByID)
		products.PUT("/:id", router.handler.HandleUpdateProductByID)
	}
}