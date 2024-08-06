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
		products.POST("", middleware.AuthMiddleware(router.p, entity.RoleUser), router.handler.HandleCreateProduct)
		products.GET("/search", middleware.AuthMiddleware(router.p), router.handler.HandleGetAllProducts)
		products.GET("", router.handler.HandleGetAllProducts)
		products.GET("/:id", middleware.AuthMiddleware(router.p), router.handler.HandleGetProductByID)
		products.DELETE("/:id", middleware.AuthMiddleware(router.p, entity.RoleAdmin), router.handler.HandleDeleteProductByID)
		products.PUT("/:id", middleware.AuthMiddleware(router.p), router.handler.HandleUpdateProductByID)
		products.GET("/report", middleware.AuthMiddleware(router.p), router.handler.HandleGetReport)
	}
}