package categories

import (
	"pm/domain/entity"
)

type CategoryRepository interface {
	Create(*entity.Category) error
	Update(*entity.Category) (*entity.Category, error)
	GetCategoryByID(int64) (*entity.Category, error)
	GetAllCategories(*entity.CategoryFilter, *entity.Pagination) ([]entity.Category, error)
	DeleteCategory(*entity.Category) error
}