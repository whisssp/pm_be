package order_items

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"pm/domain/entity"
	orderItems "pm/domain/repository/order_items"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/persistences/base"
)

type OrderItemRepository struct {
	db         *gorm.DB
	c          *gin.Context
	p          *base.Persistence
	parentSpan trace.Span
}

func (o OrderItemRepository) CreateNewOrderItems(items []entity.OrderItem) error {
	if err := o.db.Debug().Create(&items).Error; err != nil {
		return payload.ErrDB(err)
	}
	return nil
}

func (o OrderItemRepository) UpdateOrderItems(items []entity.OrderItem) ([]entity.OrderItem, error) {
	db := o.db.Begin()
	if err := db.Error; err != nil {
		return nil, payload.ErrDB(err)
	}

	for k, _ := range items {
		if err := db.Debug().Updates(&items[k]).Error; err != nil {
			db.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, payload.ErrEntityNotFound("order_items", err)
			}
			return nil, payload.ErrDB(err)
		}
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return nil, payload.ErrDB(err)
	}

	return items, nil
}

func (o OrderItemRepository) GetOrderItemByID(id int64) (*entity.OrderItem, error) {
	var orderItem entity.OrderItem
	if err := o.db.Debug().Where("id = ?", id).First(&orderItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, payload.ErrEntityNotFound("order_items", err)
		}
		return nil, payload.ErrDB(err)
	}
	return &orderItem, nil
}

func (o OrderItemRepository) GetAllOrderItems() ([]entity.OrderItem, error) {
	var orderItems []entity.OrderItem
	if err := o.db.Debug().Find(&orderItems).Error; err != nil {
		return nil, payload.ErrDB(err)
	}
	return orderItems, nil
}

func (o OrderItemRepository) DeleteOrderItemByID(id int64) error {
	db := o.db.Begin()
	if err := db.Error; err != nil {
		return payload.ErrDB(err)
	}

	orderItem, err := o.GetOrderItemByID(id)
	if err != nil {
		return err
	}

	if err := db.Model(&entity.OrderItem{}).Debug().Delete(&orderItem).Error; err != nil {
		return payload.ErrDB(err)
	}

	if err := db.Commit().Error; err != nil {
		return payload.ErrDB(err)
	}

	return nil
}

func NewOrderItemRepository(db *gorm.DB, c *gin.Context, p *base.Persistence) orderItems.OrderItemRepository {
	return OrderItemRepository{
		db: db,
		c:  c,
		p:  p,
	}
}