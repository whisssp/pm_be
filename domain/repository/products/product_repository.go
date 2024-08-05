package products

import (
	"go.opentelemetry.io/otel/trace"
	"pm/domain/entity"
)

type ProductRepository interface {
	Create(trace.Span, *entity.Product) error
	Update(trace.Span, *entity.Product) (*entity.Product, error)
	GetProductByID(trace.Span, int64) (*entity.Product, error)
	GetAllProducts(trace.Span, *entity.ProductFilter, *entity.Pagination) ([]entity.Product, error)
	DeleteProduct(trace.Span, *entity.Product) error
	GetProductByOrderItem(trace.Span, ...entity.OrderItem) ([]entity.Product, error)
	UpdateMultiProduct(trace.Span, ...entity.Product) ([]entity.Product, error)
	IsAvailableStockByOrderItems(trace.Span, ...entity.OrderItem) ([]entity.Product, error)
}