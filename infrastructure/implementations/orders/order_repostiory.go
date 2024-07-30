package orders

import (
	"fmt"
	"gorm.io/gorm"
	"pm/domain/entity"
	"pm/domain/repository"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/implementations/products"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) repository.OrderRepository {
	return OrderRepository{db}
}

func (o OrderRepository) Create(order *entity.Order) error {
	tx := o.db.Begin()
	productRepo := products.NewProductRepository(tx)
	prods, err := productRepo.GetStockByProductIDs(order.OrderItems...)
	//prodPointers := make([]*entity.Product, len(prods))
	if err != nil {
		tx.Rollback()
		return payload.ErrDB(err)
	}

	for index, v := range prods {
		v.Stock = v.Stock - int64(order.OrderItems[index].Quantity)
		if v.Stock < 0 {
			tx.Rollback()
			return payload.ErrInvalidRequest(fmt.Errorf("the product %v is out of stock", v.ID))
		}
		prods[index] = v
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	go func() {
		fmt.Println("Running goroutine update products")
		prodRepo := products.NewProductRepository(o.db)
		prods, err = prodRepo.UpdateMultiProduct(prods...)
		if err != nil {
			fmt.Printf("error updating quantity product by goroutine")
		}
		fmt.Printf("\nafter updating %v", prods)
	}()
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
