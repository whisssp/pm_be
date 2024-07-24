package repository

import (
	"pm/domain/entity"
	"pm/infrastructure/controllers/payload"
)

type ProductRepository interface {
	Create(*entity.Product) error
	Update(*entity.Product) (*entity.Product, error)
	GetProductByID(id int64) (*entity.Product, error)
	GetAllProducts(filter *payload.ProductFilter, pagination *payload.Pagination) ([]entity.Product, error)
	DeleteProduct(product *entity.Product) error
}