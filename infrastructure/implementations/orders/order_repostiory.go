package orders

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"pm/domain/entity"
	"pm/domain/repository"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/implementations/products"
	"pm/infrastructure/persistences/base"
)

type OrderRepository struct {
	db *gorm.DB
	p  *base.Persistence
	c  *gin.Context
}

func NewOrderRepository(c *gin.Context, p *base.Persistence, db *gorm.DB) repository.OrderRepository {
	return OrderRepository{db, p, c}
}

func (o OrderRepository) Create(order *entity.Order) error {
	span := o.p.Logger.Start(o.c, "CREATE_ORDER_DATABASE", o.p.Logger.SetContextWithSpanFunc())
	defer span.End()
	tx := o.db.Begin()
	o.p.Logger.Info("CREATE_ORDER_DATABASE", map[string]interface{}{"data": order})

	productRepo := products.NewProductRepository(o.c, o.p, tx)
	// Get product by product id in order item
	prods, err := productRepo.GetProductByOrderItem(order.OrderItems...)
	if err != nil {
		o.p.Logger.Error("CREATE_ORDER_DATABASE", map[string]interface{}{"message": err.Error()})
		tx.Rollback()
		return err
	}

	for index, v := range prods {
		v.Stock = v.Stock - int64(order.OrderItems[index].Quantity)
		if v.Stock < 0 {
			errS := fmt.Errorf("the product %v is out of stock", v.ID)
			o.p.Logger.Error("CREATE_ORDER_DATABASE_FAILED", map[string]interface{}{"message": errS})
			tx.Rollback()
			return payload.ErrInvalidRequest(errS)
		}
		prods[index] = v
	}

	if err := tx.Create(&order).Error; err != nil {
		o.p.Logger.Error("CREATE_ORDER_DATABASE_FAILED", map[string]interface{}{"message": err.Error()})
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		o.p.Logger.Error("CREATE_ORDER_DATABASE_FAILED", map[string]interface{}{"message": err.Error()})
		return err
	}

	go func() {
		o.p.Logger.Info("GOROUTINE_UPDATE_QUANTITY_PRODUCT", map[string]interface{}{"products": prods})
		prodRepo := products.NewProductRepository(o.c, o.p, o.db)
		prods, err = prodRepo.UpdateMultiProduct(prods...)
		if err != nil {
			zap.S().Errorw("GOROUTINE_UPDATE_QUANTITY_PRODUCT_FAILED", map[string]interface{}{"products": prods, "message": err.Error()})
		}
		zap.S().Infow("GOROUTINE_UPDATE_QUANTITY_PRODUCT_SUCCESSFULLY", map[string]interface{}{"products": prods})
	}()

	o.p.Logger.Info("CREATE_ORDER_DATABASE_SUCCESSFULLY", map[string]interface{}{"data": order})
	return nil
}

func (o OrderRepository) Update(order *entity.Order) (*entity.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o OrderRepository) GetOrderByID(id int64) (*entity.Order, error) {
	var order entity.Order
	//err := o.db.Model(&order).Where("order.id = ?", id).Association("OrderItems").Find(&order.OrderItems)
	err := o.db.Preload("OrderItems").First(&order, id).Error
	if err != nil {
		return nil, err
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