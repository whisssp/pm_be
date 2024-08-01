package jobs

import (
	"fmt"
	"go.uber.org/zap"
	"pm/domain/entity"
	"pm/infrastructure/persistences/base"
	"pm/utils"
)

const redisProductKey = "products"

func LoadProductToRedis(p *base.Persistence) {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("error trying to initialize logger")
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	// Use the logger
	sugar.Debugw("GO_ROUTINE_LOAD_PRODUCT_TO_REDIS")
	products := make([]entity.Product, 0)
	err = p.GormDB.Model(&entity.Product{}).Find(&products).Error
	if err != nil {
		sugar.Errorw("ERROR_LOAD_PRODUCT_TO_REDIS", map[string]interface{}{"message": err.Error()})
	}
	if err != nil {
		sugar.Errorw("ERROR_LOAD_PRODUCT_TO_REDIS", map[string]interface{}{"message": err.Error()})
	}

	err = utils.RedisSetHashGenericKeySlice(redisProductKey, products, entity.GetID, p.Redis.KeyExpirationTime)
	if err != nil {
		sugar.Errorw("ERROR_LOAD_PRODUCT_TO_REDIS", map[string]interface{}{"message": err.Error()})
		return
	}
}