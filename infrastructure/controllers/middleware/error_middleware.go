package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pm/infrastructure/controllers/payload"
)

func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if customErr, ok := err.(*payload.AppError); ok {
				c.JSON(customErr.StatusCode, customErr)
				return
			} else {
				c.JSON(http.StatusInternalServerError, payload.ErrInternal(err))
				return
			}
		}
	}
}