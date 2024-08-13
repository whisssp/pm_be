package products

import (
	"github.com/gin-gonic/gin"
	"pm/domain/entity"
)

type ProductRepository interface {
	Create(*entity.Product) error
	Update(*entity.Product) (*entity.Product, error)
	GetProductByID(int64) (*entity.Product, error)
	GetAllProducts(*entity.ProductFilter, *entity.Pagination) ([]entity.Product, error)
	DeleteProduct(*entity.Product) error
	GetProductByOrderItem(...entity.OrderItem) ([]entity.Product, error)
	UpdateMultiProduct(...entity.Product) ([]entity.Product, error)
	IsAvailableStockByOrderItems(*gin.Context, ...entity.OrderItem) ([]entity.Product, error)
}
