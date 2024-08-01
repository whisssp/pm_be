package application

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"math"
	"pm/domain/entity"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/implementations/orders"
	"pm/infrastructure/jobs"
	"pm/infrastructure/mapper"
	"pm/infrastructure/persistences/base"
	"pm/utils"
	"strconv"
)

const orderEntity string = "orders"

type OrderUsecase interface {
	CreateOrder(*gin.Context, *payload.CreateOrderRequest) error
	GetAllOrders(filter *entity.OrderFilter, pagination *entity.Pagination) (*payload.ListOrderResponses, error)
	GetOrderByID(*gin.Context, int64) (*payload.OrderResponse, error)
	DeleteOrderByID(id int64) error
	UpdateOrderByID(id int64, updatePayload payload.UpdateOrderRequest) (*payload.OrderResponse, error)
}

type orderUsecase struct {
	p *base.Persistence
}

func NewOrderUsecase(p *base.Persistence) OrderUsecase {
	return orderUsecase{p}
}

func (o orderUsecase) CreateOrder(c *gin.Context, reqPayload *payload.CreateOrderRequest) error {
	span := o.p.Logger.Start(c, "CREATE_ORDER_USECASES", o.p.Logger.SetContextWithSpanFunc())
	defer span.End()
	order := mapper.CreateOrderPayloadToOrder(reqPayload)
	prods := make([]entity.Product, len(order.OrderItems))

	// check product is still available stock must be greater or equal 0 after stock - quantity from order item
	for i, e := range order.OrderItems {
		var p entity.Product
		utils.RedisGetHashGenericKey(redisHashKey, strconv.Itoa(int(e.ProductID)), &p)
		var isAvailable bool = false
		if p.ID == 0 {
			continue
		}
		isAvailable = (p.Stock - int64(e.Quantity)) >= 0
		if !isAvailable {
			errP := fmt.Errorf("product: %v is not available, please check again", e.ProductID)
			o.p.Logger.Error("CREATE_ORDER_FAILED", map[string]interface{}{"message": errP.Error()})
			return payload.ErrInvalidRequest(errP)
		}
		p.Stock = p.Stock - int64(e.Quantity)
		prods[i] = p
	}

	orderRepo := orders.NewOrderRepository(c, o.p, o.p.GormDB)
	o.p.Logger.Info("CREATE_ORDER", map[string]interface{}{"order": order})
	if err := orderRepo.Create(&order); err != nil {
		o.p.Logger.Error("CREATE_ORDER_FAILED", map[string]interface{}{"message": err.Error()})
		if err, ok := err.(*payload.AppError); ok {
			return err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return payload.ErrEntityNotFound(entityName, err)
		}
		return payload.ErrDB(err)
	}

	go func() {
		logger, err := zap.NewProduction()
		if err != nil {
			fmt.Println("error trying to initialize logger")
		}
		defer logger.Sync()
		sugar := logger.Sugar()

		sugar.Debugw("GOROUTINE_UPDATE_PRODUCT_STOCK", map[string]interface{}{"products": prods})
		errMap := make(map[string]interface{})
		for _, e := range prods {
			err := utils.RedisSetHashGenericKey(redisHashKey, strconv.Itoa(int(e.ID)), e, o.p.Redis.KeyExpirationTime)
			if err != nil {
				errMap[strconv.Itoa(int(e.ID))] = err
				zap.S().Errorw("GOROUTINE_UPDATE_PRODUCT_STOCK", map[string]interface{}{"products": prods})
				go jobs.LoadProductToRedis(o.p)
			}
		}
		if len(errMap) > 0 {
			sugar.Errorw("GOROUTINE_UPDATE_PRODUCT_STOCK_FAILED", errMap)
		}
	}()
	o.p.Logger.Info("CREATE_ORDER_SUCCESSFULLY", map[string]interface{}{"order": order})
	return nil
}

func (o orderUsecase) GetAllOrders(filter *entity.OrderFilter, pagination *entity.Pagination) (*payload.ListOrderResponses, error) {
	//TODO implement me
	panic("implement me")
}

func (o orderUsecase) GetOrderByID(c *gin.Context, id int64) (*payload.OrderResponse, error) {
	db := o.p.GormDB
	orderRepo := orders.NewOrderRepository(c, o.p, db)
	order, err := orderRepo.GetOrderByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, payload.ErrEntityNotFound(orderEntity, err)
		}
		return nil, err
	}
	orderResponse := mapper.OrderToOrderResponse(order)
	var totalPrice float64 = 0
	for _, v := range order.OrderItems {
		totalPrice += v.Price * float64(v.Quantity)
	}
	orderResponse.Total = math.Round(totalPrice*100) / 100
	return &orderResponse, nil
}

func (o orderUsecase) DeleteOrderByID(id int64) error {
	//TODO implement me
	panic("implement me")
}

func (o orderUsecase) UpdateOrderByID(id int64, updatePayload payload.UpdateOrderRequest) (*payload.OrderResponse, error) {
	//TODO implement me
	panic("implement me")
}