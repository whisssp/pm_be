package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"pm/infrastructure/controllers/payload"
	"pm/utils"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer "
	userContextKey      = "user"
	subject             = "subject"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader(authorizationHeader)
		if bearerToken == "" || !strings.Contains(bearerToken, strings.TrimSpace(bearerPrefix)) {
			c.AbortWithStatusJSON(http.StatusBadRequest, payload.ErrInvalidRequest(errors.New("missing token")))
			return
		}
		tokenStr := strings.TrimSpace(strings.Split(bearerToken, bearerPrefix)[1])
		jwtToken, err := utils.ValidateToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, payload.ErrInvalidToken(err))
			return
		}
		if !jwtToken.Valid {
			c.AbortWithStatusJSON(http.StatusBadRequest, payload.ErrInvalidToken(errors.New("invalid token")))
			return
		}
		claims := utils.GetMapClaims(jwtToken)
		c.Set(userContextKey, claims[subject])
		c.Next()
	}
}