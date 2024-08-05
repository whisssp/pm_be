package application

import (
	"fmt"
	"pm/domain/entity"
	"pm/infrastructure/implementations/products"
	"pm/infrastructure/persistences/base"
	"pm/utils"
)

const (
	redisProductKey  = "products"
	redisCategoryKey = "categories"
)

type CacheUsecase interface {
	LoadProductToRedis()
}

type cacheUsecase struct {
	p *base.Persistence
}

func (c cacheUsecase) LoadProductToRedis() {
	productRepo := products.NewProductRepository(nil, c.p, c.p.GormDB)
	prods, err := productRepo.GetAllProducts(nil, nil, nil)
	if err != nil {
		fmt.Println("\nLoadProduct/error getting products from db", err)
	}
	err = utils.RedisSetHashGenericKeySlice(redisProductKey, prods, entity.GetID, c.p.Redis.KeyExpirationTime)
	if err != nil {
		fmt.Println("error adding data", err)
		return
	}
}

func NewCacheUsecase(p *base.Persistence) CacheUsecase {
	return cacheUsecase{p}
}