package repository

import (
	"pm/domain/entity"
	"pm/infrastructure/controllers/payload"
)

type CategoryRepository interface {
	Create(*entity.Category) error
	Update(*entity.Category) (*entity.Category, error)
	GetCategoryByID(id int64) (*entity.Category, error)
	GetAllCategories(filter *payload.CategoryFilter, pagination *payload.Pagination) ([]entity.Category, error)
	DeleteCategory(*entity.Category) error
}