package orders

import (
	"github.com/gin-gonic/gin"
	"pm/domain/entity"
	"pm/infrastructure/persistences/base"
)

type OrderRepository interface {
	Create(*entity.Order) error
	Update(*entity.Order) (*entity.Order, error)
	GetOrderByID(id int64) (*entity.Order, error)
	GetAllOrders(pagination *entity.Pagination) ([]entity.Order, error)
	DeleteOrder(*entity.Order) error
	IsAvailableStockByOrderItems(*base.Persistence, *gin.Context, ...entity.OrderItem) ([]entity.Product, error)
	//GetOrdersByUserID(userID int64) ([]entity.Order, error)
}