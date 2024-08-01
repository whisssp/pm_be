package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/implementations/users"
	"pm/infrastructure/persistences/base"
	"pm/utils"
	"slices"
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

func AuthMiddleware(p *base.Persistence, roles ...string) gin.HandlerFunc {
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
		user, err := users.NewUserRepository(c, p, p.GormDB).GetUserByID(idInt)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, payload.NewUnauthorized(errors.New("unauthorized"), "Not found the user from token", "ErrInvalidClaims"))
			return
		}

		if len(roles) > 0 {
			claims := utils.JwtGetMapClaims(jwtToken)
			roleFromClaims := claims[roleKey]

			if user.Role == roleFromClaims && !slices.Contains(roles, user.Role) {
				c.AbortWithStatusJSON(http.StatusForbidden, payload.ErrPermissionDenied(errors.New("You don't have permission to access this resource")))
				return
			}
		}
		c.Set(userContextKey, id)

		c.Next()
	}
}