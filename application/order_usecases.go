package application

import (
	"errors"
	"gorm.io/gorm"
	"math"
	"pm/domain/entity"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/implementations/orders"
	"pm/infrastructure/mapper"
	"pm/infrastructure/persistences/base"
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
	orderRepo := orders.NewOrderRepository(o.p.GormDB)

	if err := orderRepo.Create(&order); err != nil {
		return payload.ErrDB(err)
	}
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