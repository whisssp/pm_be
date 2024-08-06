package routes

import (
	"github.com/gin-gonic/gin"
	"pm/infrastructure/controllers/handlers"
	"pm/infrastructure/controllers/middleware"
	"pm/infrastructure/persistences/base"
)

type MailRoutes struct {
	p       *base.Persistence
	handler *handlers.MailHandler
}

func NewMailRoutes(p *base.Persistence, handler *handlers.MailHandler) *MailRoutes {
	return &MailRoutes{
		p:       p,
		handler: handler,
	}
}

func (router *MailRoutes) RegisterRoutes(routerGroup *gin.RouterGroup) {
	mailRoute := routerGroup.Group("/mail").Use(middleware.AuthMiddleware(router.p))
	{
		mailRoute.POST("/send", router.handler.HandleSendEmail)
	}
}