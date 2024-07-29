package jobs

import (
	"fmt"
	"pm/domain/entity"
	"pm/infrastructure/implementations/products"
	"pm/infrastructure/persistences/base"
	"pm/utils"
)

const redisProductKey = "products"

func LoadProductToRedis(p *base.Persistence) {
	productGormRepo := products.NewProductRepository(p.GormDB)
	prods, err := productGormRepo.GetAllProducts(nil, nil)
	if err != nil {
		fmt.Println("error getting products from database", err)
	}

	err = utils.RedisSetHashGenericKeySlice(redisProductKey, prods, entity.GetID, p.Redis.KeyExpirationTime)
	if err != nil {
		fmt.Println("error adding data", err)
		return
	}
}