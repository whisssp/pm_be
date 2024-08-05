package users

import (
	"go.opentelemetry.io/otel/trace"
	"pm/domain/entity"
)

type UserRepository interface {
	Create(trace.Span, *entity.User) error
	Update(trace.Span, *entity.User) (*entity.User, error)
	GetAllUsers(trace.Span) ([]entity.User, error)
	GetUserByID(trace.Span, int64) (*entity.User, error)
	GetUserByRole(trace.Span, entity.UserRole) (*entity.User, error)
	Delete(trace.Span, *entity.User) error
	GetUserByEmail(trace.Span, string) (*entity.User, error)
}