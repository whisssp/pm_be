package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/implementations/users"
	"pm/infrastructure/persistences/base"
	"pm/utils"
	"strconv"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer "
	userContextKey      = "user"
	subject             = "subject"
	roleKey             = "role"
)

func AuthMiddleware(p *base.Persistence, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader(authorizationHeader)
		if bearerToken == "" || !strings.Contains(bearerToken, strings.TrimSpace(bearerPrefix)) {
			c.AbortWithStatusJSON(http.StatusBadRequest, payload.ErrInvalidRequest(errors.New("missing token")))
			return
		}
		tokenStr := strings.TrimSpace(strings.Split(bearerToken, bearerPrefix)[1])
		jwtToken, err := utils.JwtValidateToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, payload.ErrInvalidToken(err))
			return
		}
		if !jwtToken.Valid {
			c.AbortWithStatusJSON(http.StatusBadRequest, payload.ErrInvalidToken(errors.New("invalid token")))
			return
		}
		id := utils.JwtGetSubject(jwtToken)
		if id == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, payload.NewUnauthorized(errors.New("unauthorized"), "invalid token", "ErrInvalidToken"))
			return
		}
		idInt, _ := strconv.ParseInt(id.(string), 10, 64)
		user, err := users.NewUserRepository(p).GetUserByID(idInt)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, payload.NewUnauthorized(errors.New("unauthorized"), "not found the user from token", "ErrInvalidClaims"))
			return
		}
		claims := utils.JwtGetMapClaims(jwtToken)
		roleFromClaims := claims[roleKey]
		if roleFromClaims != user.Role {
			c.AbortWithStatusJSON(http.StatusForbidden, payload.NewUnauthorized(errors.New("unauthorized"), "You need full permission to use this resource", "ErrNoPermission"))
			return
		}

		c.Set(userContextKey, id)
		c.Next()
	}
}