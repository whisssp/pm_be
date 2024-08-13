package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"pm/domain/entity"
	"pm/infrastructure/config"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/implementations/user_roles"
	"pm/infrastructure/persistences/base"
	"pm/infrastructure/persistences/base/logger"
	"strconv"
	"time"
)

var persistence *base.Persistence = nil

const (
	userIDKey    string = "userID"
	userRoleKey         = "role"
	expiredAtKey        = "expiredAt"
	subjectKey          = "subject"
)

var jwtSecretKey string = ""
var jwtTokenExpiration time.Duration = 10 * time.Minute

func InitJwtHelper(p *base.Persistence, jwtConfig config.JwtConfig) {
	jwtSecretKey = jwtConfig.SecretKey
	jwtTokenExpiration = jwtConfig.TokenExpiration
	persistence = p
}

func JwtGenerateJwtToken(c *gin.Context, p *base.Persistence, user *entity.User) (string, error) {
	ctx1, newlogger := logger.GetLogger().Start(c, "GENERATE_TOKEN_JWT")
	defer newlogger.End()
	newlogger.Info("STARTING_GENERATE_TOKEN", map[string]interface{}{"data": user})

	expiration := time.Now().Add(10 * 24 * time.Hour).Unix()
	//expiration := time.Now().Add(jwtTokenExpiration).Unix()
	arrBytesKey := []byte(jwtSecretKey)

	claims := jwt.MapClaims{
		subjectKey: strconv.Itoa(int(user.ID)),
		//userIDKey:    user.ID,
		userRoleKey:  user.RoleID,
		expiredAtKey: expiration,
	}
	newlogger.Info("GENERATE_TOKEN: CLAIMS", map[string]interface{}{
		"claims": claims,
	})

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(arrBytesKey)
	if err != nil {
		newlogger.Error("GENERATE_TOKEN: FAILED", map[string]interface{}{"error": err.Error()})
		return "", payload.ErrGenerateToken(err)
	}

	newlogger.Info("GENERATE_TOKEN: SUCCESSFULLY", map[string]interface{}{
		"token": token,
	})

	userRoleRepo := user_roles.NewUserRoleRepository(p.GormDB, p, ctx1)
	_, err = userRoleRepo.GetUserRoleByID(user.RoleID)
	if err != nil {
		newlogger.Error("GENERATE_TOKEN_JWT_FAILED", map[string]interface{}{"error": err.Error()})
		return "", err
	}
	return token, nil
}

func JwtGetMapClaims(token *jwt.Token) jwt.MapClaims {
	return token.Claims.(jwt.MapClaims)
}

func JwtGetSubject(token *jwt.Token) interface{} {
	if token != nil {
		claims := JwtGetMapClaims(token)
		subject := claims[subjectKey]
		if subject == nil {
			fmt.Println("error getting subject from jwt token")
			return nil
		}
		return subject
	}
	return nil
}

func JwtValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		expFloat := JwtGetMapClaims(t)[expiredAtKey].(float64)
		if time.Now().Compare(time.Unix(int64(expFloat), 0)) >= 1 {
			return nil, fmt.Errorf("token expired")
		}
		return []byte(jwtSecretKey), nil
	})
}
