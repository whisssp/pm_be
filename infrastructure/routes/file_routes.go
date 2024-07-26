package routes

import (
	"github.com/gin-gonic/gin"
	"pm/infrastructure/controllers/handlers"
)

type FileRoutes struct {
	handler *handlers.FileHandler
}

func NewFileRoutes(handler *handlers.FileHandler) *FileRoutes {
	return &FileRoutes{handler}
}

func (router *FileRoutes) RegisterRoutes(routerGroup *gin.RouterGroup) {
	files := routerGroup.Group("/files")
	{
		files.POST("/upload/images", router.handler.HandleUploadImage)
	}
}