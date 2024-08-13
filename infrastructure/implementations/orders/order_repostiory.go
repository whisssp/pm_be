package orders

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"pm/domain/entity"
	"pm/domain/repository/orders"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/persistences/base"
	"pm/infrastructure/persistences/base/logger"
)

type OrderRepository struct {
	db *gorm.DB
	p  *base.Persistence
	c  *gin.Context
}

func NewOrderRepository(c *gin.Context, p *base.Persistence, db *gorm.DB) orders.OrderRepository {
	return OrderRepository{db, p, c}
}

func (o OrderRepository) IsAvailableStockByOrderItems(persistence *base.Persistence, c *gin.Context, orderItems ...entity.OrderItem) ([]entity.Product, error) {

	persistence.Logger.Info("CHECK_STOCK", map[string]interface{}{"data": orderItems})

	ps := make([]entity.Product, 0)
	for _, o := range orderItems {
		var p entity.Product
		if err := persistence.GormDB.Model(&entity.Product{}).Where("id = ?", o.ProductID).First(&p).Error; err != nil {
			persistence.Logger.Error("CHECK_STOCK_FAILED", map[string]interface{}{"message": err.Error()})
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, payload.ErrEntityNotFound("products", err)
			}
			return nil, payload.ErrDB(err)
		}
		if p.Stock-int64(o.Quantity) < 0 {
			err := fmt.Errorf("the product %v is out of stock", o.ProductID)
			persistence.Logger.Error("CHECK_STOCK_FAILED", map[string]interface{}{"message": err.Error()})
			return nil, payload.ErrInvalidRequest(err)
		}
		ps = append(ps, p)
	}

	persistence.Logger.Info("CHECK_STOCK_SUCCESSFULLY", map[string]interface{}{"products": ps})
	return ps, nil
}

func (o OrderRepository) Create(ctx *gin.Context, order *entity.Order) error {
	_, newlogger := logger.GetLogger().Start(ctx, "create ode repository")
	defer newlogger.End()
	tx := o.db.Begin()

	newlogger.Info("Starting: create order repository", map[string]interface{}{"order": order})
	if err := tx.Create(&order).Error; err != nil {
		newlogger.Error("create order repository: error", map[string]interface{}{"error": err.Error()})
		tx.Rollback()
		return payload.ErrDB(err)
	}

	if err := tx.Commit().Error; err != nil {
		newlogger.Error("create order repository: error commit transaction", map[string]interface{}{"error": err.Error()})
		return payload.ErrDB(err)
	}

	newlogger.Info("create order repository: successfully", map[string]interface{}{"order": order})
	return nil
}

func (o OrderRepository) Update(order *entity.Order) (*entity.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o OrderRepository) GetOrderByID(id int64) (*entity.Order, error) {
	//logg := o.p.Logger
	//span := logg.Start(o.c, "GET_ORDER_BY_ID: DATABASE")
	//defer span.End()
	//o.p.Logger.

	var order entity.Order
	//err := o.db.Model(&order).Where("order.id = ?", id).Association("OrderItems").Find(&order.OrderItems)
	err := o.db.Preload("OrderItems").First(&order, id).Error
	if err != nil {
		return nil, payload.ErrDB(err)
	}
	return &order, nil
}

func (o OrderRepository) GetAllOrders(pagination *entity.Pagination) ([]entity.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o OrderRepository) DeleteOrder(order *entity.Order) error {
	//TODO implement me
	panic("implement me")
}
