package categories

import (
	"go.opentelemetry.io/otel/trace"
	"pm/domain/entity"
)

type CategoryRepository interface {
	Create(trace.Span, *entity.Category) error
	Update(trace.Span, *entity.Category) (*entity.Category, error)
	GetCategoryByID(trace.Span, int64) (*entity.Category, error)
	GetAllCategories(trace.Span, *entity.CategoryFilter, *entity.Pagination) ([]entity.Category, error)
	DeleteCategory(trace.Span, *entity.Category) error
}