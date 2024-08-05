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

func AuthMiddleware(p *base.Persistence, roles ...int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		span := p.Logger.Start(c, "AUTH_MIDDLEWARE", p.Logger.SetContextWithSpanFunc())
		defer span.End()
		p.Logger.Info("AUTH_MIDDLEWARE", map[string]interface{}{})

		bearerToken := c.GetHeader(authorizationHeader)
		if bearerToken == "" || !strings.Contains(bearerToken, strings.TrimSpace(bearerPrefix)) {
			errT := errors.New("missing token")
			p.Logger.Error("AUTHENTICATION_FAILED", map[string]interface{}{"error": errT})
			c.AbortWithStatusJSON(http.StatusBadRequest, payload.ErrInvalidRequest(errT))
			return
		}
		tokenStr := strings.TrimSpace(strings.Split(bearerToken, bearerPrefix)[1])
		jwtToken, err := utils.JwtValidateToken(tokenStr)
		if err != nil {
			p.Logger.Error("AUTHENTICATION_FAILED", map[string]interface{}{"error": err.Error()})
			c.AbortWithStatusJSON(http.StatusBadRequest, payload.ErrInvalidToken(err))
			return
		}
		if !jwtToken.Valid {
			errV := errors.New("invalid token")
			p.Logger.Error("AUTHENTICATION_FAILED", map[string]interface{}{"error": errV.Error()})
			c.AbortWithStatusJSON(http.StatusBadRequest, payload.ErrInvalidToken(errV))
			return
		}
		id := utils.JwtGetSubject(jwtToken)
		if id == nil {
			errV := payload.NewUnauthorized(errors.New("unauthorized"), "invalid token", "ErrInvalidToken")
			p.Logger.Error("AUTHENTICATION_FAILED", map[string]interface{}{"error": errV.Error()})
			c.AbortWithStatusJSON(http.StatusUnauthorized, errV)
			return
		}
		idInt, _ := strconv.ParseInt(id.(string), 10, 64)
		user, err := users.NewUserRepository(c, p, p.GormDB).GetUserByID(span, idInt)
		if err != nil {
			errU := payload.NewUnauthorized(errors.New("unauthorized"), "Not found the user from token", "ErrInvalidClaims")
			p.Logger.Error("AUTHENTICATION_FAILED", map[string]interface{}{"error": errU.Error()})
			c.AbortWithStatusJSON(http.StatusUnauthorized, errU)
			return
		}

		if len(roles) > 0 {
			claims := utils.JwtGetMapClaims(jwtToken)
			roleFromClaims, _ := strconv.ParseInt(claims[roleKey].(string), 10, 64)

			if user.RoleID == roleFromClaims && !slices.Contains(roles, roleFromClaims) {
				errP := payload.ErrPermissionDenied(errors.New("You don't have permission to access this resource"))
				p.Logger.Error("AUTHORIZATION_FAILED", map[string]interface{}{"error": errP.Error()})
				c.AbortWithStatusJSON(http.StatusForbidden, errP)
				return
			}
		}
		c.Set(userContextKey, id)

		p.Logger.Info("AUTH_MIDDLEWARE_SUCCESSFULLY", map[string]interface{}{})

		c.Next()
	}
}