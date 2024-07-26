package routes

import (
	"github.com/gin-gonic/gin"
	"pm/infrastructure/controllers/handlers"
	"pm/infrastructure/persistences/base"
)

type FileRoutes struct {
	handler *handlers.FileHandler
	p       *base.Persistence
}

func NewFileRoutes(p *base.Persistence, handler *handlers.FileHandler) *FileRoutes {
	return &FileRoutes{handler, p}
}

func (router *FileRoutes) RegisterRoutes(routerGroup *gin.RouterGroup) {
	files := routerGroup.Group("/files")
	{
		files.POST("/upload/images", router.handler.HandleUploadImage)
	}
}