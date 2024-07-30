package application

import (
	"errors"
	"fmt"
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
	CreateOrder(reqPayload *payload.CreateOrderRequest) error
	GetAllOrders(filter *entity.OrderFilter, pagination *entity.Pagination) (*payload.ListOrderResponses, error)
	GetOrderByID(id int64) (*payload.OrderResponse, error)
	DeleteOrderByID(id int64) error
	UpdateOrderByID(id int64, updatePayload payload.UpdateOrderRequest) (*payload.OrderResponse, error)
}

type orderUsecase struct {
	p *base.Persistence
}

func NewOrderUsecase(p *base.Persistence) OrderUsecase {
	return orderUsecase{p}
}

func (o orderUsecase) CreateOrder(reqPayload *payload.CreateOrderRequest) error {
	order := mapper.CreateOrderPayloadToOrder(reqPayload)
	prods := make([]entity.Product, len(order.OrderItems))

	// check product is still available stock must be greater or equal 0 after stock - quantity from order item
	for i, e := range order.OrderItems {
		var p entity.Product
		utils.RedisGetHashGenericKey(redisHashKey, strconv.Itoa(int(e.ProductID)), &p)
		isAvailable := (p.Stock - int64(e.Quantity)) >= 0
		if !isAvailable {
			return payload.ErrInvalidRequest(fmt.Errorf("product: %v is not available, please check again", e.ProductID))
		}
		p.Stock = p.Stock - int64(e.Quantity)
		prods[i] = p
	}

	orderRepo := orders.NewOrderRepository(o.p.GormDB)
	if err := orderRepo.Create(&order); err != nil {
		return payload.ErrDB(err)
	}

	go func() {
		fmt.Println("setting new product on redis")
		for _, e := range prods {
			err := utils.RedisSetHashGenericKey(redisHashKey, strconv.Itoa(int(e.ID)), e, o.p.Redis.KeyExpirationTime)
			if err != nil {
				fmt.Printf("error when updating product stock to redis: %v", prods)
				go jobs.LoadProductToRedis(o.p)
			}
		}
	}()
	return nil
}

func (o orderUsecase) GetAllOrders(filter *entity.OrderFilter, pagination *entity.Pagination) (*payload.ListOrderResponses, error) {
	//TODO implement me
	panic("implement me")
}

func (o orderUsecase) GetOrderByID(id int64) (*payload.OrderResponse, error) {
	db := o.p.GormDB
	orderRepo := orders.NewOrderRepository(db)
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
