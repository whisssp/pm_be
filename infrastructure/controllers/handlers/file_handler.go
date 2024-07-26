package handlers

import (
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"net/http"
	"pm/application"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/persistences/base"
	"pm/utils"
)

type FileHandler struct {
	p           *base.Persistence
	fileUsecase application.FileUsecase
}

func NewFileHandler(p *base.Persistence) *FileHandler {
	fileUsecase := application.NewFileUsecase(p)
	return &FileHandler{p, fileUsecase}
}

// HandleUploadImage Upload file 			godoc
// @Summary 			Upload a file
// @Description			Upload a file to get link get media
// Tag					File
// @Param				file formData file true "file"
// @Success				200		{object} payload.AppResponse{}
// @Failure      		400  	{object} payload.AppError{}
// @Failure 			500 	{object} payload.AppError{}
// @Router				/files/upload/image [post]
func (h FileHandler) HandleUploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequest(err))
		return
	}
	fileName := utils.SanitizeFileName(header.Filename)
	mime, err := mimetype.DetectReader(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, payload.ErrDetectFileType(err))
		return
	}
	if _, err := file.Seek(0, 0); err != nil {
		c.JSON(http.StatusInternalServerError, payload.ErrDetectFileType(err))
		return
	}
	defer file.Close()
	contentType := mime.String()
	r, err := h.fileUsecase.UploadImage(fileName, file, contentType)
	if err != nil {
		c.JSON(http.StatusBadRequest, payload.ErrUploadFile(err))
		return
	}
	c.JSON(http.StatusOK, payload.SuccessResponse(r, ""))
}