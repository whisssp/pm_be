package application

import (
	"fmt"
	"github.com/robfig/cron"
	"math"
	"pm/domain/entity"
	"pm/domain/repository"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/implementations/cache"
	"pm/infrastructure/implementations/products"
	"pm/infrastructure/persistences/base"
	"pm/utils"
	"strconv"
	"strings"
	"time"
)

const redisHashKey = "product"

type ProductUsecase struct {
	productGormRepo repository.ProductRepository
	redisCacheRepo  repository.RedisCacheRepository
}

func NewProductUsecase(p *base.Persistence) *ProductUsecase {
	productGormRepo := products.NewProductRepository(p.GormDB)
	redisCacheRepo := cache.NewRedisCacheRepository(p.RedisDB, p.Ctx)
	utils.InitCacheHelper(p)

	c := cron.New()
	err := c.AddFunc("0/5 * * * *", func() {
		LoadProduct(productGormRepo, redisCacheRepo)
	})

	if err != nil {
		fmt.Println("Error adding cron job:", err)
	}

	go func() {
		c.Run()
	}()

	return &ProductUsecase{
		productGormRepo: productGormRepo,
		redisCacheRepo:  redisCacheRepo,
	}
}

func LoadProduct(productGormRepo repository.ProductRepository, redisCacheRepo repository.RedisCacheRepository) {
	prods, err := productGormRepo.GetAllProducts(nil, nil)
	if err != nil {
		fmt.Println("error getting products from database", err)
	}

	err = utils.RedisSetHashGenericKeySlice(redisHashKey, prods, entity.GetID, 10*time.Minute)
	if err != nil {
		fmt.Println("error adding data", err)
		return
	}
}

func (prodUsecase *ProductUsecase) CreateProduct(reqPayload *payload.CreateProductRequest) error {
	if err := utils.ValidateReqPayload(reqPayload); err != nil {
		fmt.Printf("error validating product: %v", err)
		return payload.ErrValidateFailed(err)
	}
	prod := prodUsecase.payloadToProduct(reqPayload)
	err := prodUsecase.productGormRepo.Create(prod)
	if err != nil {
		fmt.Printf("error creating product: %v", err)
		return payload.ErrDB(err)
	}

	err = utils.RedisSetHashGenericKey(redisHashKey, strconv.FormatInt(int64(prod.ID), 10), prod, 10*time.Minute)
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			fmt.Println("error adding product to redis", err)
			LoadProduct(prodUsecase.productGormRepo, prodUsecase.redisCacheRepo)
		}

		fmt.Println("got error on redis", err)
	}
	return nil
}

func (prodUsecase *ProductUsecase) GetAllProducts(filter *payload.ProductFilter, pagination *payload.Pagination) (*payload.ListProductResponse, error) {
	var listProdResponse *payload.ListProductResponse
	if filter == nil {
		productsMap := make(map[string]entity.Product, 0)
		utils.GetAllHashGeneric(redisHashKey, &productsMap)
		prods := make([]entity.Product, 0)
		if len(productsMap) != 0 {
			prods = prodUsecase.prodsMapToArray(productsMap)
			_prods := prods[pagination.GetLimit():]
			listProdResponse = prodUsecase.prodsToProdListResponse(_prods, &payload.Pagination{
				Limit:      pagination.GetLimit(),
				Page:       pagination.GetPage(),
				Sort:       "",
				TotalRows:  int64(len(prods)),
				TotalPages: int(math.Ceil(float64(len(prods) * 1.0 / pagination.GetLimit()))),
			})
		}
		return listProdResponse, nil
	}
	prods, err := prodUsecase.productGormRepo.GetAllProducts(filter, pagination)
	if err != nil {
		return nil, payload.ErrDB(err)
	}
	listProdResponse = prodUsecase.prodsToProdListResponse(prods, pagination)
	return listProdResponse, nil
}

func (prodUsecase *ProductUsecase) GetProductByID(id int64) (*payload.ProductResponse, error) {
	var prod entity.Product
	utils.RedisGetHashGenericKey(redisHashKey, strconv.FormatInt(int64(prod.ID), 10), &prod)
	if prod.ID == 0 {
		_prod, err := prodUsecase.productGormRepo.GetProductByID(id)
		if err != nil {
			return nil, payload.ErrEntityNotFound("products", err)
		}
		prodResponse := prodUsecase.prodToProdResponse(_prod)
		return &prodResponse, nil
	}
	prodResponse := prodUsecase.prodToProdResponse(&prod)
	return &prodResponse, nil
}

func (prodUsecase *ProductUsecase) DeleteProductByID(id int64) error {
	err := utils.RedisRemoveHashGenericKey(redisHashKey, strconv.FormatInt(int64(id), 10))
	if err != nil {
		fmt.Printf("error deleting on redis: key: %v - error: %v", id, err)
	}
	prod, err := prodUsecase.productGormRepo.GetProductByID(id)
	if err != nil {
		return payload.ErrEntityNotFound("products", err)
	}
	err = prodUsecase.productGormRepo.DeleteProduct(prod)
	if err != nil {
		return err
	}
	return nil
}

func (prodUsecase *ProductUsecase) UpdateProductByID(id int64, updatePayload *payload.UpdateProductRequest) (*payload.ProductResponse, error) {
	prod, err := prodUsecase.productGormRepo.GetProductByID(id)
	if err != nil {
		return nil, payload.ErrEntityNotFound("products", err)
	}
	updatePayload.ID = id
	prodUsecase.updateProduct(prod, updatePayload)
	_, err = prodUsecase.productGormRepo.Update(prod)
	if err != nil {
		return nil, err
	}

	err = utils.RedisSetHashGenericKey(redisHashKey, strconv.FormatInt(int64(prod.ID), 10), prod, 10*time.Minute)
	if err != nil {
		fmt.Printf("error updating product: ID: %v - error: %v", id, err)
	}
	prodResponse := prodUsecase.prodToProdResponse(prod)
	return &prodResponse, nil
}

func (prodUsecase *ProductUsecase) prodsMapToArray(prodsMap map[string]entity.Product) []entity.Product {
	arrProds := make([]entity.Product, 0)
	for _, v := range prodsMap {
		arrProds = append(arrProds, v)
	}
	return arrProds
}

func (prodUsecase *ProductUsecase) prodToProdResponse(product *entity.Product) payload.ProductResponse {
	return payload.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CategoryID:  product.CategoryID,
		Stock:       product.Stock,
		Image:       product.Image,
		UpdatedAt:   product.UpdatedAt,
		CreatedAt:   product.CreatedAt,
	}
}

func (prodUsecase *ProductUsecase) prodsToProdListResponse(products []entity.Product, pagination *payload.Pagination) *payload.ListProductResponse {
	listProdResponse := make([]payload.ProductResponse, 0)
	for _, p := range products {
		prodResponse := prodUsecase.prodToProdResponse(&p)
		listProdResponse = append(listProdResponse, prodResponse)
	}
	return &payload.ListProductResponse{
		Products:      listProdResponse,
		Limit:         pagination.Limit,
		Page:          pagination.Page,
		TotalElements: pagination.TotalRows,
		TotalPages:    pagination.TotalPages,
	}
}

func (prodUsecase *ProductUsecase) payloadToProduct(reqPayload *payload.CreateProductRequest) *entity.Product {
	return &entity.Product{
		Name:        reqPayload.Name,
		Description: reqPayload.Description,
		Price:       reqPayload.Price,
		CategoryID:  reqPayload.CategoryID,
		Stock:       reqPayload.Stock,
		Image:       reqPayload.Image,
	}
}

func (prodUsecase *ProductUsecase) updateProduct(oldProd *entity.Product, updatePayload *payload.UpdateProductRequest) {
	oldProd.Name = updatePayload.Name
	oldProd.Description = updatePayload.Description
	oldProd.Stock = updatePayload.Stock
	oldProd.Price = updatePayload.Price
	oldProd.Image = updatePayload.Image
}