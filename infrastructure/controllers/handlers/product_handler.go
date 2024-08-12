package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"pm/application"
	"pm/domain/entity"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/mapper"
	"pm/infrastructure/persistences/base"
	"strconv"
	"strings"
)

const (
	entityName string = "products"
)

type ProductHandler struct {
	p       *base.Persistence
	usecase application.ProductUsecase
}

func NewProductHandler(p *base.Persistence) *ProductHandler {
	usecase := application.NewProductUsecase(p)
	return &ProductHandler{p, usecase}
}

// HandleCreateProduct CreateProduct godoc
//
//	@Summary		Create a product
//	@Description	create a new product
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			CreateProductRequest	body		payload.CreateProductRequest	true	"create product wth create product request"
//	@Success		200						{object}	payload.AppResponse
//	@Failure		400						{object}	payload.AppError
//	@Failure		404						{object}	payload.AppError
//	@Failure		500						{object}	payload.AppError
//	@Router			/products 							[post]
func (handler *ProductHandler) HandleCreateProduct(c *gin.Context) {

	var createProdReq payload.CreateProductRequest
	if err := c.ShouldBindJSON(&createProdReq); err != nil {
		c.Error(payload.ErrInvalidRequest(err))
		handler.p.Logger.Error("CREATE_PRODUCT: FAILED", map[string]interface{}{"message": err.Error()})
		return
	}

	if err := handler.usecase.CreateProduct(c, &createProdReq); err != nil {
		c.Error(err)
		handler.p.Logger.Error("CREATE_PRODUCT: FAILED", map[string]interface{}{"message": err.Error()})
		return
	}

	handler.p.Logger.Info("CREATE_PRODUCT: SUCCESSFULLY", map[string]interface{}{})
	c.JSON(http.StatusOK, payload.SuccessResponse(nil, ""))
}

// HandleGetAllProducts GetAllProducts godoc
//
//	@Summary		Get all products
//	@Description	Get all products which is not deleted
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			limit		query		int						false	"the limit perpage"
//	@Param			page		query		int						false	"the page nummber"
//	@Param			filter		query		entity.ProductFilter	false	"filtering the data"
//	@Success		200			{object}	payload.AppResponse
//	@Failure		400			{object}	payload.AppError
//	@Failure		500			{object}	payload.AppError
//	@Router			/products 				[get]
//	@Router			/products/search 				[get]
func (handler *ProductHandler) HandleGetAllProducts(c *gin.Context) {

	var filter entity.ProductFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.Error(payload.ErrInvalidRequest(err))
		handler.p.Logger.Error("GET_ALL_PRODUCTS_FAILED", map[string]interface{}{"message": err.Error()})
		return
	}

	pagination := entity.InitPaginate()
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.Error(payload.ErrInvalidRequest(err))
		handler.p.Logger.Error("GET_ALL_PRODUCTS_FAILED", map[string]interface{}{"message": err.Error()})
		return
	}

	prods, err := handler.usecase.GetAllProducts(c, &filter, pagination)
	if err != nil {
		c.Error(err)
		handler.p.Logger.Error("GET_ALL_PRODUCTS_FAILED", map[string]interface{}{"message": err.Error()})
		return
	}

	listProdResponse := mapper.ProdsToListProdsResponse(prods, pagination)
	handler.p.Logger.Info("GET_ALL_PRODUCTS_SUCCESSFULLY", map[string]interface{}{"list_products_response": listProdResponse})
	c.JSON(http.StatusOK, payload.SuccessResponse(listProdResponse, ""))
}

// HandleGetProductByID GetProductByID godoc
//
//	@Summary		Get product by id
//	@Description	Get product by id
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			id				path		int	true	"the id of the product to return"
//	@Success		200				{object}	payload.AppResponse
//	@Failure		400				{object}	payload.AppError
//	@Failure		404				{object}	payload.AppError
//	@Failure		500				{object}	payload.AppError
//	@Router			/products/:id 				[get]
func (handler *ProductHandler) HandleGetProductByID(c *gin.Context) {

	id, _ := strconv.ParseInt(removeSlashFromParam(c.Param("id")), 10, 64)
	if id == 0 {
		err := fmt.Errorf("[id] parameter is required")
		c.Error(payload.ErrInvalidRequest(err))
		handler.p.Logger.Error("GET_PRODUCT_FAILED", map[string]interface{}{"error": err.Error()})
		return
	}

	prod, err := handler.usecase.GetProductByID(c, id)
	if err != nil {
		c.Error(err)
		handler.p.Logger.Error("GET_PRODUCT_FAILED", map[string]interface{}{"error": err.Error()})
		return
	}

	prodResponse := mapper.ProductToProductResponse(prod)
	handler.p.Logger.Info("GET_PRODUCT_SUCCESSFULLY", map[string]interface{}{"product": prodResponse})
	c.JSON(http.StatusOK, payload.SuccessResponse(prodResponse, ""))
}

// HandleDeleteProductByID DeleteProductByID godoc
//
//	@Summary		Delete product by id
//	@Description	Delete product by id
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			id				path		int	true	"the id of product to delete"
//	@Success		200				{object}	payload.AppResponse
//	@Failure		400				{object}	payload.AppError
//	@Failure		404				{object}	payload.AppError
//	@Failure		500				{object}	payload.AppError
//	@Router			/products/:id 				[delete]
func (handler *ProductHandler) HandleDeleteProductByID(c *gin.Context) {

	id, _ := strconv.ParseInt(removeSlashFromParam(c.Param("id")), 10, 64)
	if id == 0 {
		err := fmt.Errorf("[id] parameter is required")
		c.Error(payload.ErrInvalidRequest(err))
		handler.p.Logger.Error("DELETE_PRODUCT_FAILED", map[string]interface{}{"message": err.Error()})
		return
	}

	err := handler.usecase.DeleteProductByID(c, id)
	if err != nil {
		c.Error(err)
		handler.p.Logger.Error("DELETE_PRODUCT_FAILED", map[string]interface{}{"message": err.Error()})
		return
	}
	handler.p.Logger.Info("DELETE_PRODUCT_SUCCESSFULLY", map[string]interface{}{"id": id})
	c.JSON(http.StatusOK, payload.SuccessResponse(nil, ""))
}

// HandleUpdateProductByID UpdateProductByID godoc
//
//	@Summary		Update product by id
//	@Description	Update product by id
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			id				path		int	true	"the id of product to update"
//	@Param			UpdateProductRequest	body		payload.UpdateProductRequest	true	"update product with update product request"
//	@Success		200				{object}	payload.AppResponse
//	@Failure		400				{object}	payload.AppError
//	@Failure		404				{object}	payload.AppError
//	@Failure		500				{object}	payload.AppError
//	@Router			/products/:id 				[put]
func (handler *ProductHandler) HandleUpdateProductByID(c *gin.Context) {

	id, _ := strconv.ParseInt(removeSlashFromParam(c.Param("id")), 10, 64)
	if id == 0 {
		err := fmt.Errorf("[id] parameter is required")
		c.Error(payload.ErrInvalidRequest(err))
		handler.p.Logger.Error("UPDATE_PRODUCT_FAILED", map[string]interface{}{"message": err.Error()})
		return
	}

	var updateProductReq payload.UpdateProductRequest
	if err := c.ShouldBindJSON(&updateProductReq); err != nil {
		c.Error(payload.ErrInvalidRequest(err))
		handler.p.Logger.Error("UPDATE_PRODUCT_FAILED", map[string]interface{}{"message": err.Error()})
		return
	}

	prodUpdated, err := handler.usecase.UpdateProductByID(c, id, &updateProductReq)
	if err != nil {
		c.Error(err)
		handler.p.Logger.Error("UPDATE_PRODUCT_FAILED", map[string]interface{}{"message": err.Error()})
		return
	}
	prodResponse := mapper.ProductToProductResponse(prodUpdated)
	handler.p.Logger.Info("UPDATE_PRODUCT_SUCCESSFULLY", map[string]interface{}{"product_response": prodResponse})
	c.JSON(http.StatusOK, payload.SuccessResponse(prodResponse, ""))
}

func (handler *ProductHandler) HandleGetReport(c *gin.Context) {
	err := handler.usecase.Report()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, payload.SuccessResponse(nil, "dang gui"))
}

func removeSlashFromParam(param string) string {
	if strings.Contains(param, "/") {
		param = strings.Replace(param, "/", "", -1)
	}
	if strings.Contains(param, "\\") {
		param = strings.Replace(param, "\\", "", -1)
	}
	return param
}
