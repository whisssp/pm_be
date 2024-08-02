package application

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"pm/domain/entity"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/implementations/products"
	"pm/infrastructure/mapper"
	"pm/infrastructure/persistences/base"
	"pm/utils"
	"strconv"
	"strings"
)

const (
	entityName   string = "products"
	redisHashKey        = "products"
)

type ProductUsecase interface {
	CreateProduct(*gin.Context, *payload.CreateProductRequest) error
	GetAllProducts(*gin.Context, *entity.ProductFilter, *entity.Pagination) (*payload.ListProductResponses, error)
	GetProductByID(*gin.Context, int64) (*payload.ProductResponse, error)
	DeleteProductByID(*gin.Context, int64) error
	UpdateProductByID(*gin.Context, int64, *payload.UpdateProductRequest) (*payload.ProductResponse, error)
}
type productUsecase struct {
	p *base.Persistence
}

func NewProductUsecase(p *base.Persistence) ProductUsecase {
	return productUsecase{p}
}

func (p productUsecase) CreateProduct(c *gin.Context, reqPayload *payload.CreateProductRequest) error {
	span := p.p.Logger.Start(c, "CREATE_PRODUCT_USECASES", p.p.Logger.SetContextWithSpanFunc())
	defer span.End()
	p.p.Logger.Info("CREATE_PRODUCT", map[string]interface{}{"data": reqPayload})

	if err := utils.ValidateReqPayload(reqPayload); err != nil {
		p.p.Logger.Error("CREATE_PRODUCT_FAILED", map[string]interface{}{"message": err.Error()})
		return payload.ErrValidateFailed(err)
	}

	prod := mapper.PayloadToProduct(reqPayload)
	productRepo := products.NewProductRepository(c, p.p, p.p.GormDB)
	err := productRepo.Create(prod)
	if err != nil {
		p.p.Logger.Error("CREATE_PRODUCT_FAILED", map[string]interface{}{"message": err.Error()})
		return err
	}

	go func(prod *entity.Product) {

		err = utils.RedisSetHashGenericKey(redisHashKey, strconv.FormatInt(int64(prod.ID), 10), prod, p.p.Redis.KeyExpirationTime)
		if err != nil {
			if !strings.Contains(err.Error(), "not found") {
				fmt.Println("error adding product to redis", err)
				prods, err := productRepo.GetAllProducts(nil, nil)
				if err != nil {
					fmt.Println("\nerror getting product from db", err)
				}
				err = utils.RedisSetHashGenericKeySlice(redisProductKey, prods, entity.GetID, p.p.Redis.KeyExpirationTime)
			}

			fmt.Println("got error on redis", err)
		}
	}(prod)

	p.p.Logger.Info("CREATE_PRODUCT_SUCCESSFULLY", map[string]interface{}{"data": prod.ID})
	return nil
}

func (p productUsecase) GetAllProducts(c *gin.Context, filter *entity.ProductFilter, pagination *entity.Pagination) (*payload.ListProductResponses, error) {
	span := p.p.Logger.Start(c, "GET_ALL_PRODUCTS_USECASES", p.p.Logger.SetContextWithSpanFunc())
	defer span.End()
	p.p.Logger.Info("GET_ALL_PRODUCTS", map[string]interface{}{"params": struct {
		Filter     interface{} `json:"filter"`
		Pagination interface{} `json:"pagination"`
	}{
		Filter:     filter,
		Pagination: pagination,
	}})

	var listProdResponse payload.ListProductResponses
	productRepo := products.NewProductRepository(c, p.p, p.p.GormDB)
	prods := make([]entity.Product, 0)
	if filter.IsNil() == true {
		productsMap := make(map[string]entity.Product)
		utils.GetAllHashGeneric(redisHashKey, &productsMap)
		prods := make([]entity.Product, 0)
		if len(productsMap) != 0 {
			prods = prodsMapToArray(productsMap)
			_prods := prods[:pagination.GetOffset()-1]
			_prods = _prods[:pagination.GetLimit()-1]
			listProdResponse = mapper.ProdsToListProdsResponse(_prods, &entity.Pagination{
				Limit:      pagination.GetLimit(),
				Page:       pagination.GetPage(),
				Sort:       "",
				TotalRows:  int64(len(prods)),
				TotalPages: int(math.Ceil(float64(len(prods) * 1.0 / pagination.GetLimit()))),
			})
			p.p.Logger.Info("GET_ALL_PRODUCTS_SUCCESSFULLY", map[string]interface{}{"data": listProdResponse})
			return &listProdResponse, nil
		}
	}
	prods, err := productRepo.GetAllProducts(filter, pagination)
	if err != nil {
		p.p.Logger.Info("GET_ALL_PRODUCTS_FAILED", map[string]interface{}{"message": err.Error()})
		return nil, err
	}

	listProdResponse = mapper.ProdsToListProdsResponse(prods, pagination)
	p.p.Logger.Info("GET_ALL_PRODUCTS_SUCCESSFULLY", map[string]interface{}{"data": listProdResponse})
	return &listProdResponse, nil
}

func (p productUsecase) GetProductByID(c *gin.Context, id int64) (*payload.ProductResponse, error) {
	span := p.p.Logger.Start(c, "GET_PRODUCT_BY_ID_USECASES", p.p.Logger.SetContextWithSpanFunc())
	defer span.End()
	p.p.Logger.Info("GET_PRODUCT", map[string]interface{}{"data": id})

	var prod entity.Product
	utils.RedisGetHashGenericKey(redisHashKey, strconv.FormatInt(int64(prod.ID), 10), &prod)
	productRepo := products.NewProductRepository(c, p.p, p.p.GormDB)
	// **set a go routine to log error from redis w zap

	if prod.ID != 0 {
		//_prod, err := productRepo.GetProductByID(id)
		//if err != nil {
		//	return nil, payload.ErrEntityNotFound(entityName, err)
		//}
		prodResponse := mapper.ProductToProductResponse(&prod)
		p.p.Logger.Info("GET_PRODUCT_SUCCESSFULLY", map[string]interface{}{"data": prodResponse})
		return &prodResponse, nil
	}

	prodPointer, err := productRepo.GetProductByID(id)
	if err != nil {
		p.p.Logger.Error("GET_PRODUCT_FAILED", map[string]interface{}{"message": err.Error()})
		return nil, err
	}
	prodResponse := mapper.ProductToProductResponse(prodPointer)
	p.p.Logger.Info("GET_PRODUCT_SUCCESSFULLY", map[string]interface{}{"data": prodResponse})
	return &prodResponse, nil
}

func (p productUsecase) DeleteProductByID(c *gin.Context, id int64) error {
	span := p.p.Logger.Start(c, "DELETE_PRODUCT_BY_ID_USECASES", p.p.Logger.SetContextWithSpanFunc())
	defer span.End()
	p.p.Logger.Info("DELETE_PRODUCT", map[string]interface{}{"data": id})

	err := utils.RedisRemoveHashGenericKey(redisHashKey, strconv.FormatInt(int64(id), 10))
	if err != nil {
		// log by zap
		fmt.Printf("error deleting on redis: key: %v - error: %v", id, err)
	}
	productRepo := products.NewProductRepository(c, p.p, p.p.GormDB)
	prod, err := productRepo.GetProductByID(id)
	if err != nil {
		p.p.Logger.Info("DELETE_PRODUCT_FAILED", map[string]interface{}{"message": err.Error()})
		return payload.ErrEntityNotFound(entityName, err)
	}
	err = productRepo.DeleteProduct(prod)
	if err != nil {
		p.p.Logger.Info("DELETE_PRODUCT_FAILED", map[string]interface{}{"message": err.Error()})
		return err
	}

	p.p.Logger.Info("DELETE_PRODUCT_SUCCESSFULLY", map[string]interface{}{})
	return nil
}

func (p productUsecase) UpdateProductByID(c *gin.Context, id int64, updatePayload *payload.UpdateProductRequest) (*payload.ProductResponse, error) {
	span := p.p.Logger.Start(c, "UPDATE_PRODUCT_USECASES", p.p.Logger.SetContextWithSpanFunc())
	defer span.End()
	p.p.Logger.Info("UPDATE_PRODUCT", map[string]interface{}{"data": struct {
		ID      interface{}
		Payload interface{}
	}{
		ID:      id,
		Payload: updatePayload,
	}})

	productRepo := products.NewProductRepository(c, p.p, p.p.GormDB)
	prod, err := productRepo.GetProductByID(id)
	if err != nil {
		p.p.Logger.Error("UPDATE_PRODUCT_FAILED", map[string]interface{}{"message": err.Error()})
		return nil, payload.ErrEntityNotFound(entityName, err)
	}
	updatePayload.ID = id
	mapper.UpdateProduct(prod, updatePayload)
	_, err = productRepo.Update(prod)
	if err != nil {
		p.p.Logger.Error("UPDATE_PRODUCT_FAILED", map[string]interface{}{"message": err.Error()})
		return nil, payload.ErrCannotUpdateEntity(entityName, err)
	}

	err = utils.RedisSetHashGenericKey(redisHashKey, strconv.FormatInt(int64(prod.ID), 10), prod, p.p.Redis.KeyExpirationTime)
	if err != nil {
		fmt.Printf("error updating product: ID: %v - error: %v", id, err)
	}
	prodResponse := mapper.ProductToProductResponse(prod)
	p.p.Logger.Error("UPDATE_PRODUCT_SUCCESSFULLY", map[string]interface{}{"data": prodResponse})
	return &prodResponse, nil
}

func prodsMapToArray(prodsMap map[string]entity.Product) []entity.Product {
	arrProds := make([]entity.Product, 0)
	for _, v := range prodsMap {
		arrProds = append(arrProds, v)
	}
	return arrProds
}