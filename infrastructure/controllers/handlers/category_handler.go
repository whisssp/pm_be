package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"pm/application"
	"pm/domain/entity"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/persistences/base"
	"strconv"
)

type CategoryHandler struct {
	p               *base.Persistence
	categoryUsecase application.CategoryUsecase
}

func NewCategoryHandler(p *base.Persistence) *CategoryHandler {
	categoryUsecase := application.NewCategoryUsecase(p)

	return &CategoryHandler{p, categoryUsecase}
}

// HandleCreateCategory CreateCategory godoc
//
//	@Summary		Create a category
//	@Description	create a new category
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Param			CreateCategoryRequest	body		payload.CreateCategoryRequest{}	true	"create category with create category request"
//	@Success		200						{object}	payload.AppResponse
//	@Failure		400						{object}	payload.AppError
//	@Failure		404						{object}	payload.AppError
//	@Failure		500						{object}	payload.AppError
//	@Router			/categories 							[post]
func (h CategoryHandler) HandleCreateCategory(c *gin.Context) {
	var categoryReq payload.CreateCategoryRequest
	if err := c.ShouldBindJSON(&categoryReq); err != nil {
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequest(err))
		return
	}

	if err := h.categoryUsecase.CreateCategory(&categoryReq); err != nil {
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequest(err))
		return
	}

	c.JSON(http.StatusOK, payload.SuccessResponse(nil, ""))
}

// HandleGetAllCategories GetAllCategories godoc
//
//	@Summary		Get all categories
//	@Description	Get all categoies which is not deleted
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Param			limit		query		int						false	"the limit perpage"
//	@Param			page		query		int						false	"the page nummber"
//	@Param			filter		query		entity.CategoryFilter	false	"filtering the data"
//	@Success		200			{object}	payload.AppResponse
//	@Failure		400			{object}	payload.AppError
//	@Failure		500			{object}	payload.AppError
//	@Router			/categories 				[get]
func (h CategoryHandler) HandleGetAllCategories(c *gin.Context) {
	var categoryFilter entity.CategoryFilter
	var pagination entity.Pagination

	if err := c.ShouldBindQuery(&categoryFilter); err != nil {
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequest(err))
		return
	}

	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequest(err))
		return
	}

	categories, err := h.categoryUsecase.GetAllCategories(&categoryFilter, &pagination)
	if err != nil {
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequest(err))
		return
	}
	c.JSON(http.StatusOK, payload.SuccessResponse(categories, ""))
}

// HandleGetCategoryByID GetCategoryByID godoc
//
//	@Summary		Get category by id
//	@Description	Get category by id
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Param			id				path		int	true	"the id of the category to return"
//	@Success		200				{object}	payload.AppResponse
//	@Failure		400				{object}	payload.AppError
//	@Failure		404				{object}	payload.AppError
//	@Failure		500				{object}	payload.AppError
//	@Router			/categories/:id 				[get]
func (h CategoryHandler) HandleGetCategoryByID(c *gin.Context) {
	id, _ := strconv.ParseInt(removeSlashFromParam(c.Param("id")), 10, 64)
	if id == 0 {
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequest(fmt.Errorf("id must be a string of numbers")))
		return
	}
	prod, err := h.categoryUsecase.GetCategoryByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, payload.ErrEntityNotFound("categories", err))
			return
		}
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequest(err))
		return
	}
	c.JSON(http.StatusOK, payload.SuccessResponse(prod, ""))
}

// HandleDeleteCategoryByID DeleteCategoryByID godoc
//
//	@Summary		Delete category by id
//	@Description	Delete category by id
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Param			id				path		int	true	"the id of category to delete"
//	@Success		200				{object}	payload.AppResponse
//	@Failure		400				{object}	payload.AppError
//	@Failure		404				{object}	payload.AppError
//	@Failure		500				{object}	payload.AppError
//	@Router			/categories/:id 				[delete]
func (h CategoryHandler) HandleDeleteCategoryByID(c *gin.Context) {
	id, _ := strconv.ParseInt(removeSlashFromParam(c.Param("id")), 10, 64)
	if id == 0 {
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequest(fmt.Errorf("id must be a string of numbers")))
		return
	}
	err := h.categoryUsecase.DeleteCategoryByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, payload.ErrEntityNotFound("categories", err))
			return
		}
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequest(err))
		return
	}
	c.JSON(http.StatusOK, payload.SuccessResponse(nil, ""))
}

// HandleUpdateCategoryByID UpdateCategoryByID godoc
//
//	@Summary		Update category by id
//	@Description	Update category by id
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Param			id				path		int	true	"the id of category to update"
//	@Param			UpdateCategoryRequest	body		payload.UpdateCategoryRequest	true	"update cateogory with update category request"
//	@Success		200				{object}	payload.AppResponse
//	@Failure		400				{object}	payload.AppError
//	@Failure		404				{object}	payload.AppError
//	@Failure		500				{object}	payload.AppError
//	@Router			/categories/:id 				[put]
func (h CategoryHandler) HandleUpdateCategoryByID(c *gin.Context) {
	var updatePayload payload.UpdateCategoryRequest
	id, _ := strconv.ParseInt(removeSlashFromParam(c.Param("id")), 10, 64)
	if id == 0 {
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequest(fmt.Errorf("id must be a string of numbers")))
		return
	}
	if err := c.ShouldBindJSON(&updatePayload); err != nil {
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequest(err))
		return
	}
	categoryUpdated, err := h.categoryUsecase.UpdateCategoryByID(id, updatePayload)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, payload.ErrEntityNotFound("categories", err))
			return
		}
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequest(err))
		return
	}
	c.JSON(http.StatusOK, payload.SuccessResponse(categoryUpdated, ""))
}