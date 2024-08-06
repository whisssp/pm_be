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
	span := o.p.Logger.Start(o.c, "CREATE_ORDER_ITEMS", o.p.Logger.UseGivenSpan(o.parentSpan))
	defer span.End()

	if err := o.db.Debug().Create(&items).Error; err != nil {
		o.p.Logger.Error("CREATE_ORDER_ITEMS: ERROR", map[string]interface{}{"error": err.Error(), "order_items": items})
		return payload.ErrDB(err)
	}
	return nil
}

func (o OrderItemRepository) UpdateOrderItems(items []entity.OrderItem) ([]entity.OrderItem, error) {
	span := o.p.Logger.Start(o.c, "UPDATE_ORDER_ITEMS: REPO", o.p.Logger.UseGivenSpan(o.parentSpan))
	defer span.End()
	db := o.db.Begin()
	if err := db.Error; err != nil {
		o.p.Logger.Error("UPDATE_ORDER_ITEMS: ERROR BEGIN TRANSACTION", map[string]interface{}{"error": err.Error()}, o.p.Logger.UseGivenSpan(span))
		return nil, payload.ErrDB(err)
	}

	for k, _ := range items {
		if err := db.Debug().Updates(&items[k]).Error; err != nil {
			db.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				o.p.Logger.Error("UPDATE_ORDER_ITEMS: ERROR ORDER_ITEM NOT FOUND", map[string]interface{}{"error": err.Error(), "order_item_id": items[k].Model.ID}, o.p.Logger.UseGivenSpan(span))
				return nil, payload.ErrEntityNotFound("order_items", err)
			}
			o.p.Logger.Error("UPDATE_ORDER_ITEMS: ERROR", map[string]interface{}{"error": err.Error()}, o.p.Logger.UseGivenSpan(span))
			return nil, payload.ErrDB(err)
		}
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		o.p.Logger.Error("UPDATE_ORDER_ITEMS: ERROR COMMIT TRANSACTION", map[string]interface{}{"error": err.Error()}, o.p.Logger.UseGivenSpan(span))
		return nil, payload.ErrDB(err)
	}

	return items, nil
}

func (o OrderItemRepository) GetOrderItemByID(id int64) (*entity.OrderItem, error) {
	span := o.p.Logger.Start(o.c, "GET_ORDER_ITEM: REPO", o.p.Logger.UseGivenSpan(o.parentSpan))
	defer span.End()
	var orderItem entity.OrderItem
	if err := o.db.Debug().Where("id = ?", id).First(&orderItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			o.p.Logger.Error("UPDATE_ORDER_ITEMS: ERROR ORDER_ITEM NOT FOUND", map[string]interface{}{"error": err.Error(), "order_item_id": id}, o.p.Logger.UseGivenSpan(span))
			return nil, payload.ErrEntityNotFound("order_items", err)
		}
		o.p.Logger.Error("GET_ORDER_ITEM: ERROR DB", map[string]interface{}{"error": err.Error()}, o.p.Logger.UseGivenSpan(span))
		return nil, payload.ErrDB(err)
	}
	return &orderItem, nil
}

func (o OrderItemRepository) GetAllOrderItems() ([]entity.OrderItem, error) {
	span := o.p.Logger.Start(o.c, "GET_ALL_ORDER_ITEM: REPO", o.p.Logger.UseGivenSpan(o.parentSpan))
	defer span.End()
	var orderItems []entity.OrderItem
	if err := o.db.Debug().Find(&orderItems).Error; err != nil {
		o.p.Logger.Error("GET_ALL_ORDER_ITEM: ERROR DB", map[string]interface{}{"error": err.Error()}, o.p.Logger.UseGivenSpan(span))
		return nil, payload.ErrDB(err)
	}
	return orderItems, nil
}

func (o OrderItemRepository) DeleteOrderItemByID(id int64) error {
	span := o.p.Logger.Start(o.c, "DELETE_ORDER_ITEM: REPO", o.p.Logger.UseGivenSpan(o.parentSpan))
	defer span.End()
	db := o.db.Begin()
	if err := db.Error; err != nil {
		o.p.Logger.Error("DELETE_ORDER_ITEM: ERROR BEGIN TRANSACTION", map[string]interface{}{"error": err.Error()}, o.p.Logger.UseGivenSpan(span))
		return payload.ErrDB(err)
	}

	orderItem, err := o.GetOrderItemByID(id)
	if err != nil {
		o.p.Logger.Error("DELETE_ORDER_ITEM: ERROR ORDER_ITEM NOT FOUND", map[string]interface{}{"error": err.Error()}, o.p.Logger.UseGivenSpan(span))
		return err
	}

	if err := db.Model(&entity.OrderItem{}).Debug().Delete(&orderItem).Error; err != nil {
		o.p.Logger.Error("DELETE_ORDER_ITEM: ERROR", map[string]interface{}{"error": err.Error()}, o.p.Logger.UseGivenSpan(span))
		return payload.ErrDB(err)
	}

	if err := db.Commit().Error; err != nil {
		o.p.Logger.Error("DELETE_ORDER_ITEM: ERROR COMMIT TRANSACTION", map[string]interface{}{"error": err.Error()}, o.p.Logger.UseGivenSpan(span))
		return payload.ErrDB(err)
	}

	return nil
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