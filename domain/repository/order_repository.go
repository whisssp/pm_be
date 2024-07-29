package repository

import "pm/domain/entity"

type OrderRepository interface {
	//Create(*entity.Order) error

	Create(*entity.Order) error
	Update(*entity.Order) (*entity.Order, error)
	GetOrderByID(id int64) (*entity.Order, error)
	GetAllOrders(pagination *entity.Pagination) ([]entity.Order, error)
	DeleteOrder(*entity.Order) error
	//GetOrdersByUserID(userID int64) ([]entity.Order, error)
	//Get
}