package application

import (
	"fmt"
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

const redisHashKey = "products"

type ProductUsecase interface {
	CreateProduct(reqPayload *payload.CreateProductRequest) error
	GetAllProducts(filter *payload.ProductFilter, pagination *payload.Pagination) (*payload.ListProductResponse, error)
	GetProductByID(id int64) (*payload.ProductResponse, error)
	DeleteProductByID(id int64) error
	UpdateProductByID(id int64, updatePayload *payload.UpdateProductRequest) (*payload.ProductResponse, error)
}
type productUsecase struct {
	p *base.Persistence
}

func NewProductUsecase(p *base.Persistence) ProductUsecase {
	return productUsecase{p}
}

func (p productUsecase) CreateProduct(reqPayload *payload.CreateProductRequest) error {
	if err := utils.ValidateReqPayload(reqPayload); err != nil {
		fmt.Printf("error validating product: %v", err)
		return payload.ErrValidateFailed(err)
	}

	prod := mapper.PayloadToProduct(reqPayload)
	productRepo := products.NewProductRepository(p.p)
	err := productRepo.Create(prod)
	if err != nil {
		fmt.Printf("error creating product: %v", err)
		return payload.ErrDB(err)
	}

	err = utils.RedisSetHashGenericKey(redisHashKey, strconv.FormatInt(int64(prod.ID), 10), prod, p.p.RedisExpirationTime)
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			fmt.Println("error adding product to redis", err)
			prods, err := productRepo.GetAllProducts(nil, nil)
			if err != nil {
				fmt.Println("\nerror getting product from db", err)
			}
			err = utils.RedisSetHashGenericKeySlice(redisProductKey, prods, entity.GetID, p.p.RedisExpirationTime)
		}

		fmt.Println("got error on redis", err)
	}
	return nil
}

func (p productUsecase) GetAllProducts(filter *payload.ProductFilter, pagination *payload.Pagination) (*payload.ListProductResponse, error) {
	var listProdResponse payload.ListProductResponse
	productRepo := products.NewProductRepository(p.p)
	prods := make([]entity.Product, 0)
	if filter.IsNil() == true {
		productsMap := make(map[string]entity.Product)
		utils.GetAllHashGeneric(redisHashKey, &productsMap)
		prods := make([]entity.Product, 0)
		if len(productsMap) != 0 {
			prods = prodsMapToArray(productsMap)
			_prods := prods[:pagination.GetOffset()-1]
			_prods = _prods[:pagination.GetLimit()-1]
			listProdResponse = mapper.ProdsToListProdsResponse(_prods, &payload.Pagination{
				Limit:      pagination.GetLimit(),
				Page:       pagination.GetPage(),
				Sort:       "",
				TotalRows:  int64(len(prods)),
				TotalPages: int(math.Ceil(float64(len(prods) * 1.0 / pagination.GetLimit()))),
			})
			return &listProdResponse, nil
		}
	}
	prods, err := productRepo.GetAllProducts(filter, pagination)
	if err != nil {
		return nil, payload.ErrDB(err)
	}

	listProdResponse = mapper.ProdsToListProdsResponse(prods, pagination)
	return &listProdResponse, nil
}

func (p productUsecase) GetProductByID(id int64) (*payload.ProductResponse, error) {
	var prod entity.Product
	utils.RedisGetHashGenericKey(redisHashKey, strconv.FormatInt(int64(prod.ID), 10), &prod)
	productRepo := products.NewProductRepository(p.p)
	if prod.ID != 0 {
		_prod, err := productRepo.GetProductByID(id)
		if err != nil {
			return nil, payload.ErrEntityNotFound("products", err)
		}
		prodResponse := mapper.ProductToProductResponse(_prod)
		return &prodResponse, nil
	}

	prodPointer, err := productRepo.GetProductByID(id)
	if err != nil {
		return nil, payload.ErrEntityNotFound("products", err)
	}
	prodResponse := mapper.ProductToProductResponse(prodPointer)
	return &prodResponse, nil
}

func (p productUsecase) DeleteProductByID(id int64) error {
	err := utils.RedisRemoveHashGenericKey(redisHashKey, strconv.FormatInt(int64(id), 10))
	if err != nil {
		fmt.Printf("error deleting on redis: key: %v - error: %v", id, err)
	}
	productRepo := products.NewProductRepository(p.p)
	prod, err := productRepo.GetProductByID(id)
	if err != nil {
		return payload.ErrEntityNotFound("products", err)
	}
	err = productRepo.DeleteProduct(prod)
	if err != nil {
		return err
	}
	return nil
}

func (p productUsecase) UpdateProductByID(id int64, updatePayload *payload.UpdateProductRequest) (*payload.ProductResponse, error) {
	productRepo := products.NewProductRepository(p.p)
	prod, err := productRepo.GetProductByID(id)
	if err != nil {
		return nil, payload.ErrEntityNotFound("products", err)
	}
	updatePayload.ID = id
	mapper.UpdateProduct(prod, updatePayload)
	_, err = productRepo.Update(prod)
	if err != nil {
		return nil, err
	}

	err = utils.RedisSetHashGenericKey(redisHashKey, strconv.FormatInt(int64(prod.ID), 10), prod, p.p.RedisExpirationTime)
	if err != nil {
		fmt.Printf("error updating product: ID: %v - error: %v", id, err)
	}
	prodResponse := mapper.ProductToProductResponse(prod)
	return &prodResponse, nil
}

func prodsMapToArray(prodsMap map[string]entity.Product) []entity.Product {
	arrProds := make([]entity.Product, 0)
	for _, v := range prodsMap {
		arrProds = append(arrProds, v)
	}
	return arrProds
}