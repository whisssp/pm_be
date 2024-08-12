package application

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"pm/domain/entity"
	"pm/infrastructure/controllers/payload"
	orderItems "pm/infrastructure/implementations/order_items"
	"pm/infrastructure/implementations/products"
	"pm/infrastructure/jobs"
	"pm/infrastructure/mapper"
	"pm/infrastructure/persistences/base"
	"pm/utils"
	"strconv"
)

type OrderItemUsecase interface {
	CreateNewOrderItem(*gin.Context, []entity.OrderItem) error
	UpdateOrderItem(*gin.Context, []entity.OrderItem) ([]payload.OrderItemResponse, error)
	GetOrderItemByID(*gin.Context, int64) (*payload.OrderItemResponse, error)
	DeleteOrderItemByID(*gin.Context, int64) error
	GetAllOrderItems(*gin.Context) (*payload.ListOrderItemResponses, error)
}

type orderItemUsecase struct {
	p *base.Persistence
}

func (o orderItemUsecase) CreateNewOrderItem(c *gin.Context, items []entity.OrderItem) error {

	for _, item := range items {
		if err := utils.ValidateReqPayload(item); err != nil {
			return payload.ErrInvalidRequest(err)
		}
	}

	prods := make([]entity.Product, 0)

	// check product stock on redis
	for i, e := range items {
		var p entity.Product
		utils.RedisGetHashGenericKey(redisHashKey, strconv.Itoa(int(e.ProductID)), &p)
		var isAvailable bool = false
		if p.ID == 0 {
			continue
		}
		isAvailable = (p.Stock - int64(e.Quantity)) >= 0
		if !isAvailable {
			err := fmt.Errorf("product: %v is not available, please check again", e.ProductID)
			zap.S().Error("Error", err)
			prods = make([]entity.Product, 0)
			break
		}
		p.Stock = p.Stock - int64(e.Quantity)
		prods[i] = p
	}

	if len(prods) == 0 {
		var err error = nil
		productRepo := products.NewProductRepository(c, o.p, o.p.GormDB)
		prods, err = productRepo.IsAvailableStockByOrderItems(nil, items...)
		if err != nil {
			return err
		}
	}

	oiRepo := orderItems.NewOrderItemRepository(o.p.GormDB, c, o.p)
	if err := oiRepo.CreateNewOrderItems(items); err != nil {
		return err
	}

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
	go func(gProds []entity.Product) {
		logger, err := zap.NewProduction()
		if err != nil {
			fmt.Println("error trying to initialize logger")
		}
		defer logger.Sync()
		sugar := logger.Sugar()

		sugar.Infow("GOROUTINE_UPDATE_PRODUCT_QUANTITY", map[string]interface{}{"data": gProds})
		prodRepo := products.NewProductRepository(c, o.p, o.p.GormDB)
		prods, err := prodRepo.UpdateMultiProduct(prods...)
		if err != nil {
			sugar.Errorw("GOROUTINE_UPDATE_QUANTITY_PRODUCT_FAILED", map[string]interface{}{"products": prods, "message": err.Error()})
		}
	}(prods)

	return nil
}

func (o orderItemUsecase) UpdateOrderItem(c *gin.Context, items []entity.OrderItem) ([]payload.OrderItemResponse, error) {

	for _, item := range items {
		if err := utils.ValidateReqPayload(item); err != nil {
			return nil, payload.ErrInvalidRequest(err)
		}

	}

	oiRepo := orderItems.NewOrderItemRepository(o.p.GormDB, c, o.p)
	orderItems, err := oiRepo.UpdateOrderItems(items)
	if err != nil {
		return nil, err
	}

	prods := make([]entity.Product, 0)

	// check product stock on redis
	for i, e := range items {
		var p entity.Product
		utils.RedisGetHashGenericKey(redisHashKey, strconv.Itoa(int(e.ProductID)), &p)
		var isAvailable bool = false
		if p.ID == 0 {
			continue
		}
		isAvailable = (p.Stock - int64(e.Quantity)) >= 0
		if !isAvailable {
			err := fmt.Errorf("product: %v is not available, please check again", e.ProductID)
			zap.S().Errorw("error", err)
			prods = make([]entity.Product, 0)
			break
		}
		p.Stock = p.Stock - int64(e.Quantity)
		prods[i] = p
	}

	if len(prods) == 0 {
		var err error = nil
		productRepo := products.NewProductRepository(c, o.p, o.p.GormDB)
		prods, err = productRepo.IsAvailableStockByOrderItems(nil, items...)
		if err != nil {
			return nil, err
		}
	}

	orderItemResponses := make([]payload.OrderItemResponse, 0)
	for _, item := range orderItems {
		oir := mapper.OrderItemToOrderItemResponse(&item)
		orderItemResponses = append(orderItemResponses, oir)
	}

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
	go func(gProds []entity.Product) {
		logger, err := zap.NewProduction()
		if err != nil {
			fmt.Println("error trying to initialize logger")
		}
		defer logger.Sync()
		sugar := logger.Sugar()

		sugar.Infow("GOROUTINE_UPDATE_PRODUCT_QUANTITY", map[string]interface{}{"data": gProds})
		prodRepo := products.NewProductRepository(c, o.p, o.p.GormDB)
		prods, err := prodRepo.UpdateMultiProduct(prods...)
		if err != nil {
			sugar.Errorw("GOROUTINE_UPDATE_QUANTITY_PRODUCT_FAILED", map[string]interface{}{"products": prods, "message": err.Error()})
		}
	}(prods)

	return orderItemResponses, nil
}

func (o orderItemUsecase) GetOrderItemByID(c *gin.Context, id int64) (*payload.OrderItemResponse, error) {

	oiRepo := orderItems.NewOrderItemRepository(o.p.GormDB, c, o.p)
	orderItem, err := oiRepo.GetOrderItemByID(id)
	if err != nil {
		return nil, err
	}

	oir := mapper.OrderItemToOrderItemResponse(orderItem)

	return &oir, nil
}

func (o orderItemUsecase) DeleteOrderItemByID(c *gin.Context, id int64) error {

	oiRepo := orderItems.NewOrderItemRepository(o.p.GormDB, c, o.p)
	err := oiRepo.DeleteOrderItemByID(id)
	if err != nil {
		return err
	}

	return nil
}

func (o orderItemUsecase) GetAllOrderItems(c *gin.Context) (*payload.ListOrderItemResponses, error) {

	oiRepo := orderItems.NewOrderItemRepository(o.p.GormDB, c, o.p)
	items, err := oiRepo.GetAllOrderItems()
	if err != nil {
		return nil, err
	}

	orderItemResponses := make([]payload.OrderItemResponse, 0)
	for _, item := range items {
		oir := mapper.OrderItemToOrderItemResponse(&item)
		orderItemResponses = append(orderItemResponses, oir)
	}
	response := payload.ListOrderItemResponses{Orders: orderItemResponses}

	return &response, nil
}

func NewOrderItemUsecase(p *base.Persistence) OrderItemUsecase {
	return orderItemUsecase{p: p}
}