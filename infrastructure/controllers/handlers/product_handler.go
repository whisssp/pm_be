package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"pm/application"
	"pm/domain/entity"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/persistences/base"
	"strconv"
	"strings"
)

type ProductHandler struct {
	p       *base.Persistence
	usecase application.ProductUsecase
}

func NewProductHandler(p *base.Persistence) *ProductHandler {
	usecase := application.NewProductUsecase(p)
	return &ProductHandler{p, usecase}
}

// CreateProduct godoc
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
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequest(err))
		return
	}

	if err := handler.usecase.CreateProduct(&createProdReq); err != nil {
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequest(err))
		return
	}

	c.JSON(http.StatusOK, payload.SuccessResponse(nil, ""))
}

// GetAllProducts godoc
//
//	@Summary		Get all products
//	@Description	Get all products which is not deleted
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			limit		query		int						false	"the limit perpage"
//	@Param			page		query		int						false	"the page nummber"
//	@Param			filter		query		payload.ProductFilter	false	"filtering the data"
//	@Success		200			{object}	payload.AppResponse
//	@Failure		400			{object}	payload.AppError
//	@Failure		500			{object}	payload.AppError
//	@Router			/products 				[get]
//	@Router			/products/search 				[get]
func (handler *ProductHandler) HandleGetAllProducts(c *gin.Context) {
	var filter entity.ProductFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequest(err))
		return
	}

	pagination := entity.InitPaginate()
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequest(err))
		return
	}

	prods, err := handler.usecase.GetAllProducts(&filter, pagination)
	if err != nil {
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequest(err))
		return
	}
	c.JSON(http.StatusOK, payload.SuccessResponse(prods, ""))
}

// GetProductByID godoc
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
		c.JSON(http.StatusBadRequest, payload.ErrParamRequired(fmt.Errorf("[id] parameter is required")))
		return
	}

	prod, err := handler.usecase.GetProductByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, payload.ErrEntityNotFound("products", err))
			return
		}
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequest(err))
		return
	}
	c.JSON(http.StatusOK, payload.SuccessResponse(prod, ""))
}

// DeleteProductByID godoc
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
		c.JSON(http.StatusBadRequest, payload.ErrParamRequired(fmt.Errorf("[id] parameter is required")))
		return
	}

	err := handler.usecase.DeleteProductByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, payload.ErrEntityNotFound("products", err))
			return
		}
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequest(err))
		return
	}
	c.JSON(http.StatusOK, payload.SuccessResponse(nil, ""))
}

// UpdateProductByID godoc
//
//	@Summary		Update product by id
//	@Description	Update product by id
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			id				path		int	true	"the id of product to update"
//	@Success		200				{object}	payload.AppResponse
//	@Failure		400				{object}	payload.AppError
//	@Failure		404				{object}	payload.AppError
//	@Failure		500				{object}	payload.AppError
//	@Router			/products/:id 				[put]
func (handler *ProductHandler) HandleUpdateProductByID(c *gin.Context) {
	id, _ := strconv.ParseInt(removeSlashFromParam(c.Param("id")), 10, 64)
	if id == 0 {
		c.JSON(http.StatusBadRequest, payload.ErrParamRequired(fmt.Errorf("[id] parameter is required")))
		return
	}

	var updateProductReq payload.UpdateProductRequest
	if err := c.ShouldBindJSON(&updateProductReq); err != nil {
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequest(err))
		return
	}

	prodUpdated, err := handler.usecase.UpdateProductByID(id, &updateProductReq)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, payload.ErrEntityNotFound("products", err))
			return
		}
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequest(err))
		return
	}
	c.JSON(http.StatusOK, payload.SuccessResponse(prodUpdated, ""))
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