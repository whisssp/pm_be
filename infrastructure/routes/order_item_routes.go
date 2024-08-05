package routes

import (
	"github.com/gin-gonic/gin"
	"pm/infrastructure/controllers/handlers"
	"pm/infrastructure/controllers/middleware"
	"pm/infrastructure/persistences/base"
)

type OrderItemRoutes struct {
	p       *base.Persistence
	handler *handlers.OrderItemHandler
}

func NewOrderItemRoutes(p *base.Persistence, handler *handlers.OrderItemHandler) *OrderItemRoutes {
	return &OrderItemRoutes{
		p:       p,
		handler: handler,
	}
}

func (router *OrderItemRoutes) RegisterRoutes(routerGroup *gin.RouterGroup) {
	orderItems := routerGroup.Group("/order-items").Use(middleware.AuthMiddleware(router.p))
	{
		orderItems.GET("/:id", router.handler.HandleGetOrderItemByID)
		orderItems.DELETE("/:id", router.handler.HandleDeleteOrderItemByID)
		orderItems.PUT("/:id", router.handler.HandleUpdateOrderItemByID)
		orderItems.GET("/:id", router.handler.HandleGetOrderItemByID)
		orderItems.GET("", router.handler.HandleGetAllOrderItems)
	}
}