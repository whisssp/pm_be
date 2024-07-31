package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"pm/domain/entity"
	"pm/infrastructure/config"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/persistences/base"
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

func JwtGenerateJwtToken(user *entity.User) (string, error) {
	expiration := time.Now().Add(10 * 24 * time.Hour).Unix()
	//expiration := time.Now().Add(jwtTokenExpiration).Unix()
	arrBytesKey := []byte(jwtSecretKey)
	claims := jwt.MapClaims{
		subjectKey: strconv.Itoa(int(user.ID)),
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