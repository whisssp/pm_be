package orders

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"pm/domain/entity"
	"pm/domain/repository"
	"pm/infrastructure/controllers/payload"
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
	span := o.p.Logger.Start(o.c, "CREATE_ORDER_DATABASE")
	defer span.End()
	o.p.Logger.Info("CREATE_ORDER_REPO", map[string]interface{}{"data": order})

	tx := o.db.Begin()

	//productRepo := products.NewProductRepository(o.c, o.p, tx)
	//// Get product by product id in order item
	//prods, err := productRepo.GetProductByOrderItem(span, order.OrderItems...)
	//if err != nil {
	//	o.p.Logger.Error("CREATE_ORDER_FAILED", map[string]interface{}{"message": err.Error()})
	//	tx.Rollback()
	//	return payload.ErrDB(err)
	//}

	//for index, v := range prods {
	//	v.Stock = v.Stock - int64(order.OrderItems[index].Quantity)
	//	if v.Stock < 0 {
	//		errS := fmt.Errorf("the product %v is out of stock", v.ID)
	//		o.p.Logger.Error("CREATE_ORDER_FAILED", map[string]interface{}{"message": errS.Error()})
	//		tx.Rollback()
	//		return payload.ErrInvalidRequest(errS)
	//	}
	//	prods[index] = v
	//}

	if err := tx.Create(&order).Error; err != nil {
		o.p.Logger.Error("CREATE_ORDER_FAILED", map[string]interface{}{"message": err.Error()})
		tx.Rollback()
		return payload.ErrDB(err)
	}

	if err := tx.Commit().Error; err != nil {
		o.p.Logger.Error("CREATE_ORDER_FAILED", map[string]interface{}{"message": err.Error()})
		return payload.ErrDB(err)
	}

	o.p.Logger.Info("CREATE_ORDER_SUCCESSFULLY", map[string]interface{}{"data": order})

	//go func(o OrderRepository) {
	//	span := o.p.Logger.Start(o.c, "GOROUTINE_UPDATE_QUANTITY_PRODUCT", o.p.Logger.SetContextWithSpanFunc())
	//	defer span.End()
	//	o.p.Logger.Info("GOROUTINE_UPDATE_QUANTITY_PRODUCT", map[string]interface{}{"products": prods})
	//	prodRepo := products.NewProductRepository(o.c, o.p, o.db)
	//	prods, err = prodRepo.UpdateMultiProduct(prods...)
	//	if err != nil {
	//		zap.S().Errorw("GOROUTINE_UPDATE_QUANTITY_PRODUCT_FAILED", map[string]interface{}{"products": prods, "message": err.Error()})
	//	}
	//	zap.S().Infow("GOROUTINE_UPDATE_QUANTITY_PRODUCT_SUCCESSFULLY", map[string]interface{}{"products": prods})
	//	o.p.Logger.Info("GOROUTINE_UPDATE_QUANTITY_PRODUCT_SUCCESSFULLY", map[string]interface{}{"products": prods})
	//}(o)

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