package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"pm/domain/entity"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/persistences/base"
	"time"
)

var persistence *base.Persistence = nil

const (
	userIDKey    string = "userID"
	userRoleKey         = "role"
	expiredAtKey        = "expiredAt"
	subjectKey          = "subject"
)

func InitJwtHelper(p *base.Persistence) {
	persistence = p
}

func GenerateJwtToken(user *entity.User) (string, error) {
	expiration := time.Now().Add(persistence.Jwt.TokenExpiration).Unix()
	arrBytesKey := []byte(persistence.Jwt.SecretKey)
	claims := jwt.MapClaims{
		subjectKey: user.ID,
		//userIDKey:    user.ID,
		userRoleKey:  user.Role,
		expiredAtKey: expiration,
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(arrBytesKey)
	if err != nil {
		return "", payload.ErrGenerateToken(err)
	}
	return token, nil
}

func GetMapClaims(token *jwt.Token) jwt.MapClaims {
	return token.Claims.(jwt.MapClaims)
}

func ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(persistence.Jwt.SecretKey), nil
	})
}