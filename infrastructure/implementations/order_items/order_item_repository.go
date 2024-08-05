package order_items

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"pm/domain/entity"
	orderItems "pm/domain/repository/order_items"
	"pm/infrastructure/persistences/base"
)

type OrderItemRepository struct {
	db         *gorm.DB
	c          *gin.Context
	p          *base.Persistence
	parentSpan trace.Span
}

func (o OrderItemRepository) CreateNewOrderItems(items []entity.OrderItem) error {
	//TODO implement me
	panic("implement me")
}

func (o OrderItemRepository) UpdateOrderItems(items []entity.OrderItem) ([]entity.OrderItem, error) {
	//TODO implement me
	panic("implement me")
}

func (o OrderItemRepository) GetOrderItemByID(i int64) (*entity.OrderItem, error) {
	//TODO implement me
	panic("implement me")
}

func (o OrderItemRepository) GetAllOrderItems() ([]entity.OrderItem, error) {
	//TODO implement me
	panic("implement me")
}

func (o OrderItemRepository) DeleteOrderItemByID(i int64) error {
	//TODO implement me
	panic("implement me")
}

func NewOrderItemRepository(db *gorm.DB, c *gin.Context, p *base.Persistence, parentSpan trace.Span) orderItems.OrderItemRepository {
	var span trace.Span = nil
	if parentSpan != nil {
		span = parentSpan
	}

	return OrderItemRepository{
		db:         db,
		c:          c,
		p:          p,
		parentSpan: span,
	}
}