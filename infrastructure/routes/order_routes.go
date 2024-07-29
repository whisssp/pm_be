package routes

import (
	"github.com/gin-gonic/gin"
	"pm/infrastructure/controllers/handlers"
	"pm/infrastructure/controllers/middleware"
	"pm/infrastructure/persistences/base"
)

type OrderRoutes struct {
	p       *base.Persistence
	handler *handlers.OrderHandler
}

func NewOrderRoutes(p *base.Persistence, handler *handlers.OrderHandler) *OrderRoutes {
	return &OrderRoutes{p, handler}
}

func (r *OrderRoutes) RegisterRoutes(routerGroup *gin.RouterGroup) {
	ordersRouter := routerGroup.Group("/orders").Use(middleware.AuthMiddleware(r.p))
	{
		ordersRouter.POST("", r.handler.HandleCreateOrder)
		ordersRouter.PUT("/:id", r.handler.HandleUpdateOrderByID)
		ordersRouter.GET("", r.handler.HandleGetAllOrders)
		ordersRouter.GET("/:id", r.handler.HandleGetOrderByID)
		ordersRouter.DELETE("/:id", r.handler.HandleDeleteOrderByID)
	}
}