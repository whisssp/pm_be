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
	span := o.p.Logger.Start(c, "CREATE_ORDER_ITEM_USECASE", o.p.Logger.SetContextWithSpanFunc())
	defer span.End()

	for _, item := range items {
		if err := utils.ValidateReqPayload(item); err != nil {
			o.p.Logger.Error("CREATE_ORDER_ITEMS: ERROR INVALID DATA", map[string]interface{}{"error": err.Error(), "order_item": item})
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
			o.p.Logger.Error("CREATE_ORDER: ERROR PRODUCT IS NOT AVAILABLE", map[string]interface{}{"error": err.Error()})
			prods = make([]entity.Product, 0)
			break
		}
		p.Stock = p.Stock - int64(e.Quantity)
		prods[i] = p
	}

	if len(prods) == 0 {
		var err error = nil
		productRepo := products.NewProductRepository(c, o.p, o.p.GormDB)
		prods, err = productRepo.IsAvailableStockByOrderItems(span, items...)
		if err != nil {
			o.p.Logger.Error("CREATE_ORDER: ERROR PRODUCT IS NOT AVAILABLE", map[string]interface{}{"error": err.Error()})
			return err
		}
	}

	oiRepo := orderItems.NewOrderItemRepository(o.p.GormDB, c, o.p, span)
	if err := oiRepo.CreateNewOrderItems(items); err != nil {
		o.p.Logger.Error("CREATE_ORDER_ITEMS: ERROR", map[string]interface{}{"error": err.Error(), "order_items": items})
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
		prods, err := prodRepo.UpdateMultiProduct(nil, prods...)
		if err != nil {
			sugar.Errorw("GOROUTINE_UPDATE_QUANTITY_PRODUCT_FAILED", map[string]interface{}{"products": prods, "message": err.Error()})
		}
	}(prods)

	o.p.Logger.Info("CREATE_ORDER_ITEMS: SUCCESSFULLY", map[string]interface{}{"order_items": items})
	return nil
}

func (o orderItemUsecase) UpdateOrderItem(c *gin.Context, items []entity.OrderItem) ([]payload.OrderItemResponse, error) {
	span := o.p.Logger.Start(c, "UPDATE_ORDER_ITEM_USECASE", o.p.Logger.SetContextWithSpanFunc())
	defer span.End()

	for _, item := range items {
		if err := utils.ValidateReqPayload(item); err != nil {
			o.p.Logger.Error("UPDATE_ORDER_ITEMS: ERROR INVALID DATA", map[string]interface{}{"error": err.Error(), "order_item": item})
			return nil, payload.ErrInvalidRequest(err)
		}

	}

	oiRepo := orderItems.NewOrderItemRepository(o.p.GormDB, c, o.p, span)
	orderItems, err := oiRepo.UpdateOrderItems(items)
	if err != nil {
		o.p.Logger.Error("UPDATE_ORDER_ITEMS: ERROR", map[string]interface{}{"error": err.Error(), "order_items": items})
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
			o.p.Logger.Error("UPDATE_ORDER_ITEM: ERROR PRODUCT IS NOT AVAILABLE", map[string]interface{}{"error": err.Error()})
			prods = make([]entity.Product, 0)
			break
		}
		p.Stock = p.Stock - int64(e.Quantity)
		prods[i] = p
	}

	if len(prods) == 0 {
		var err error = nil
		productRepo := products.NewProductRepository(c, o.p, o.p.GormDB)
		prods, err = productRepo.IsAvailableStockByOrderItems(span, items...)
		if err != nil {
			o.p.Logger.Error("UPDATE_ORDER_ITEM: ERROR PRODUCT IS NOT AVAILABLE", map[string]interface{}{"error": err.Error()})
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
		prods, err := prodRepo.UpdateMultiProduct(nil, prods...)
		if err != nil {
			sugar.Errorw("GOROUTINE_UPDATE_QUANTITY_PRODUCT_FAILED", map[string]interface{}{"products": prods, "message": err.Error()})
		}
	}(prods)

	o.p.Logger.Info("UPDATE_ORDER_ITEMS: SUCCESSFULLY", map[string]interface{}{"order_items": items})
	return orderItemResponses, nil
}

func (o orderItemUsecase) GetOrderItemByID(c *gin.Context, id int64) (*payload.OrderItemResponse, error) {
	span := o.p.Logger.Start(c, "GET_ORDER_ITEM_BY_ID_USECASE", o.p.Logger.SetContextWithSpanFunc())
	defer span.End()

	oiRepo := orderItems.NewOrderItemRepository(o.p.GormDB, c, o.p, span)
	orderItem, err := oiRepo.GetOrderItemByID(id)
	if err != nil {
		o.p.Logger.Error("GET_ORDER_ITEM_BY_ID: ERROR", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	oir := mapper.OrderItemToOrderItemResponse(orderItem)

	o.p.Logger.Info("GET_ORDER_ITEM_BY_ID: SUCCESSFULLY", map[string]interface{}{"order_item": oir})
	return &oir, nil
}

func (o orderItemUsecase) DeleteOrderItemByID(c *gin.Context, id int64) error {
	span := o.p.Logger.Start(c, "DELETE_ORDER_ITEM_BY_ID_USECASE", o.p.Logger.SetContextWithSpanFunc())
	defer span.End()

	oiRepo := orderItems.NewOrderItemRepository(o.p.GormDB, c, o.p, span)
	err := oiRepo.DeleteOrderItemByID(id)
	if err != nil {
		o.p.Logger.Error("DELETE_ORDER_ITEM_BY_ID: ERROR", map[string]interface{}{"error": err.Error()})
		return err
	}

	o.p.Logger.Info("DELETE_ORDER_ITEM_BY_ID: SUCCESSFULLY", map[string]interface{}{"deleted_id": id})
	return nil
}

func (o orderItemUsecase) GetAllOrderItems(c *gin.Context) (*payload.ListOrderItemResponses, error) {
	span := o.p.Logger.Start(c, "GET_ALL_ORDER_ITEMS_USECASE", o.p.Logger.SetContextWithSpanFunc())
	defer span.End()

	oiRepo := orderItems.NewOrderItemRepository(o.p.GormDB, c, o.p, span)
	items, err := oiRepo.GetAllOrderItems()
	if err != nil {
		o.p.Logger.Error("GET_ALL_ORDER_ITEMS: ERROR", map[string]interface{}{"error": err.Error(), "order_items": items})
		return nil, err
	}

	orderItemResponses := make([]payload.OrderItemResponse, 0)
	for _, item := range items {
		oir := mapper.OrderItemToOrderItemResponse(&item)
		orderItemResponses = append(orderItemResponses, oir)
	}
	response := payload.ListOrderItemResponses{Orders: orderItemResponses}

	o.p.Logger.Info("GET_ALL_ORDER_ITEMS: SUCCESSFULLY", map[string]interface{}{"order_items": response})
	return &response, nil
}

func NewOrderItemUsecase(p *base.Persistence) OrderItemUsecase {
	return orderItemUsecase{p: p}
}