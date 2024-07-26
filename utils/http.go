package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"pm/infrastructure/controllers/payload"
)

func HTTPSuccessResponse(ctx *gin.Context, data interface{}, message string) {
	ctx.JSON(http.StatusOK, payload.SuccessResponse(data, message))
}

func HttpErrorResponse(ctx *gin.Context, err error) {
	appErr, ok := err.(*payload.AppError)
	if ok {
		ctx.JSON(appErr.StatusCode, appErr)
		return
	}
	ctx.JSON(http.StatusBadRequest, payload.ErrInvalidRequest(err))
}

func HttpNotFoundResponse(ctx *gin.Context, err error) {
	var appErr *payload.AppError
	ok := errors.As(err, &appErr)
	if !ok {
		ctx.JSON(http.StatusNotFound, payload.ErrEntityNotFound("", err))
		return
	}
	ctx.JSON(http.StatusNotFound, appErr)
}

func HttpUnauthorizedResponse(ctx *gin.Context, err error) {

}