package application

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"pm/domain/entity"
	products2 "pm/domain/repository/products"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/implementations/orders"
	"pm/infrastructure/implementations/products"
	"pm/infrastructure/jobs"
	"pm/infrastructure/mapper"
	"pm/infrastructure/persistences/base"
	"pm/infrastructure/persistences/base/logger"
	"pm/utils"
	"strconv"
)

const orderEntity string = "orders"

type OrderUsecase interface {
	CreateOrder(*gin.Context, *payload.CreateOrderRequest) error
	GetAllOrders(filter *entity.OrderFilter, pagination *entity.Pagination) ([]entity.Order, error)
	GetOrderByID(*gin.Context, int64) (*entity.Order, error)
	DeleteOrderByID(id int64) error
	UpdateOrderByID(id int64, updatePayload payload.UpdateOrderRequest) (*entity.Order, error)
}

type orderUsecase struct {
	p *base.Persistence
}

func NewOrderUsecase(p *base.Persistence) OrderUsecase {
	return orderUsecase{p}
}

func (o orderUsecase) CreateOrder(c *gin.Context, reqPayload *payload.CreateOrderRequest) error {
	ctx, newLogger := logger.GetLogger().Start(c, "CREATE_ORDER: USECASES")
	defer newLogger.End()
	newLogger.Info("STARTING: CREATE_ORDER", map[string]interface{}{"data": reqPayload})

	order := mapper.CreateOrderPayloadToOrder(reqPayload)
	prods := make([]entity.Product, 0)
	productRepo := products.NewProductRepository(ctx, o.p, o.p.GormDB)
	orderRepo := orders.NewOrderRepository(ctx, o.p, o.p.GormDB)
	var err error

	// check product stock on redis
	for i, e := range order.OrderItems {
		var p entity.Product
		utils.RedisGetHashGenericKey(redisHashKey, strconv.Itoa(int(e.ProductID)), &p)
		var isAvailable bool = false
		if p.ID == 0 {
			continue
		}
		isAvailable = (p.Stock - int64(e.Quantity)) >= 0
		if !isAvailable {
			//err := fmt.Errorf("product: %v is not available, please check again", e.ProductID)
			//o.p.Logger.Error("CREATE_ORDER: ERROR PRODUCT IS NOT AVAILABLE", map[string]interface{}{"error": err.Error()})
			prods = make([]entity.Product, 0)
			break
		}
		p.Stock = p.Stock - int64(e.Quantity)
		prods[i] = p
	}

	if len(prods) == 0 {
		prods, err = productRepo.IsAvailableStockByOrderItems(ctx, order.OrderItems...)
		if err != nil {
			newLogger.Error("CREATE_ORDER: ERROR PRODUCT IS NOT AVAILABLE", map[string]interface{}{"error": err.Error()})
			return err
		}
	}

	newLogger.Info("CREATE_ORDER", map[string]interface{}{"order": order})
	if err := orderRepo.Create(ctx, &order); err != nil {
		newLogger.Error("CREATE_ORDER: ERROR", map[string]interface{}{"error": err.Error()})
		return err
	}
	newLogger.Info("CREATE_ORDER: SUCCESSFULLY", map[string]interface{}{"order": order})

	// this goroutine is for updating product on redis
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
				sugar.Errorw("GOROUTINE_UPDATE_PRODUCT_STOCK", map[string]interface{}{"products": prods})
				go jobs.LoadProductToRedis(o.p)
			}
		}
	}()

	// this goroutine is for updating product stock
	go func(prodRepo products2.ProductRepository, gProds []entity.Product) {
		logger, err := zap.NewProduction()
		if err != nil {
			fmt.Println("error trying to initialize logger")
		}
		defer logger.Sync()
		sugar := logger.Sugar()

		sugar.Infow("GOROUTINE_UPDATE_PRODUCT_QUANTITY", map[string]interface{}{"data": gProds})
		prods, err := prodRepo.UpdateMultiProduct(prods...)
		if err != nil {
			sugar.Errorw("GOROUTINE_UPDATE_QUANTITY_PRODUCT_FAILED", map[string]interface{}{"products": prods, "message": err.Error()})
		}
	}(productRepo, prods)
	return nil
}

func (o orderUsecase) GetAllOrders(filter *entity.OrderFilter, pagination *entity.Pagination) ([]entity.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o orderUsecase) GetOrderByID(c *gin.Context, id int64) (*entity.Order, error) {
	ctx, _ := o.p.Logger.Start(c, "GET_ORDER_BY_ID: USECASES")
	defer o.p.Logger.End()
	o.p.Logger.Info("STARTING: GET ORDER BY ID", map[string]interface{}{"id": id})

	db := o.p.GormDB
	orderRepo := orders.NewOrderRepository(ctx, o.p, db)
	order, err := orderRepo.GetOrderByID(id)
	if err != nil {
		o.p.Logger.Error("GET_ORDER: ERROR", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	o.p.Logger.Error("GET_ORDER: SUCCESSFULLY", map[string]interface{}{"order_response": order})
	return order, nil
}

func (o orderUsecase) DeleteOrderByID(id int64) error {
	//TODO implement me
	panic("implement me")
}

func (o orderUsecase) UpdateOrderByID(id int64, updatePayload payload.UpdateOrderRequest) (*entity.Order, error) {
	//TODO implement me
	panic("implement me")
}
