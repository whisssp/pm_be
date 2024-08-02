package repository

import (
	"go.opentelemetry.io/otel/trace"
	"pm/domain/entity"
)

type ProductRepository interface {
	Create(*entity.Product) error
	Update(*entity.Product) (*entity.Product, error)
	GetProductByID(id int64) (*entity.Product, error)
	GetAllProducts(filter *entity.ProductFilter, pagination *entity.Pagination) ([]entity.Product, error)
	DeleteProduct(product *entity.Product) error
	GetProductByOrderItem(trace.Span, ...entity.OrderItem) ([]entity.Product, error)
	UpdateMultiProduct(...entity.Product) ([]entity.Product, error)
	IsAvailableStockByOrderItems(...entity.OrderItem) ([]entity.Product, error)
}