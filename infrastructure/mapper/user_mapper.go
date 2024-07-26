package mapper

import (
	"pm/domain/entity"
	"pm/infrastructure/controllers/payload"
)

func UserRequestToUser(userRequest *payload.UserRequest) *entity.User {
	role := userRequest.Role.String()
	return &entity.User{
		Name:     userRequest.Name,
		Phone:    userRequest.Phone,
		Email:    userRequest.Email,
		Password: userRequest.Password,
		Role:     role,
	}
}