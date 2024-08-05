package order_items

import "pm/domain/entity"

type OrderItemRepository interface {
	CreateNewOrderItems([]entity.OrderItem) error
	UpdateOrderItems([]entity.OrderItem) ([]entity.OrderItem, error)
	GetOrderItemByID(int64) (*entity.OrderItem, error)
	GetAllOrderItems() ([]entity.OrderItem, error)
	DeleteOrderItemByID(int64) error
}