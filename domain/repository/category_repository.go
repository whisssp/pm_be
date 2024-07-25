package repository

import (
	"pm/domain/entity"
)

type CategoryRepository interface {
	Create(*entity.Category) error
	Update(*entity.Category) (*entity.Category, error)
	GetCategoryByID(id int64) (*entity.Category, error)
	GetAllCategories(filter *entity.CategoryFilter, pagination *entity.Pagination) ([]entity.Category, error)
	DeleteCategory(*entity.Category) error
}