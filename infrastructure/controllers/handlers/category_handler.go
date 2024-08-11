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
	span := h.p.Logger.Start(c, "handlers/HandleCreateCategory", h.p.Logger.SetContextWithSpanFunc())
	defer span.End()

	var categoryReq payload.CreateCategoryRequest
	if err := c.ShouldBindJSON(&categoryReq); err != nil {
		c.Error(payload.ErrInvalidRequest(err))
		h.p.Logger.Error("CREATE_CATEGORY_FAILED", map[string]interface{}{"message": err.Error()})
		return
	}

	if err := h.categoryUsecase.CreateCategory(c, &categoryReq); err != nil {
		c.Error(err)
		h.p.Logger.Error("CREATE_CATEGORY_DATABASE_FAILED", map[string]interface{}{"message": err.Error()})
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
	span := h.p.Logger.Start(c, "handlers/GetAllCategories", h.p.Logger.SetContextWithSpanFunc())
	defer span.End()

	var categoryFilter entity.CategoryFilter
	var pagination entity.Pagination

	if err := c.ShouldBindQuery(&categoryFilter); err != nil {
		c.Error(payload.ErrInvalidRequest(err))
		h.p.Logger.Error("GET_ALL_CATEGORIES_FAILED", map[string]interface{}{"message": err.Error()})
		return
	}

	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.Error(payload.ErrInvalidRequest(err))
		h.p.Logger.Error("GET_ALL_CATEGORIES_FAILED", map[string]interface{}{"message": err.Error()})
		return
	}

	categories, err := h.categoryUsecase.GetAllCategories(c, &categoryFilter, &pagination)
	if err != nil {
		c.Error(payload.ErrInvalidRequest(err))
		h.p.Logger.Error("GET_ALL_CATEGORIES_FAILED", map[string]interface{}{"message": err.Error()})
		return
	}

	listCatesResponse := mapper.CategoriesToListCategoriesResponse(categories, &pagination)
	c.JSON(http.StatusOK, payload.SuccessResponse(listCatesResponse, ""))
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
	span := h.p.Logger.Start(c, "handlers/HandleGetCategoryByID", h.p.Logger.SetContextWithSpanFunc())
	defer span.End()
	h.p.Logger.Info("GET_CATEGORY", map[string]interface{}{})

	id, _ := strconv.ParseInt(removeSlashFromParam(c.Param("id")), 10, 64)
	if id == 0 {
		err := fmt.Errorf("id must be a string of numbers")
		c.Error(payload.ErrInvalidRequest(err))
		h.p.Logger.Error("GET_CATEGORY_FAILED", map[string]interface{}{"data": err.Error()})
		return
	}
	prod, err := h.categoryUsecase.GetCategoryByID(c, id)
	if err != nil {
		c.Error(err)
		h.p.Logger.Error("GET_CATEGORY", map[string]interface{}{"data": err.Error()})
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
	span := h.p.Logger.Start(c, "handlers/HandleDeleteCategoryByID", h.p.Logger.SetContextWithSpanFunc())
	defer span.End()
	h.p.Logger.Info("DELETE_CATEGORY", map[string]interface{}{})

	id, _ := strconv.ParseInt(removeSlashFromParam(c.Param("id")), 10, 64)
	if id == 0 {
		err := fmt.Errorf("id must be a string of numbers")
		c.Error(payload.ErrInvalidRequest(err))
		h.p.Logger.Error("DELETE_CATEGORY_FAILED", map[string]interface{}{"message": err.Error()})
		return
	}
	err := h.categoryUsecase.DeleteCategoryByID(c, id)
	if err != nil {
		c.Error(err)
		h.p.Logger.Error("DELETE_CATEGORY_FAILED", map[string]interface{}{"message": err.Error()})
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
	span := h.p.Logger.Start(c, "handlers/HandleUpdateCategoryByID", h.p.Logger.SetContextWithSpanFunc())
	defer span.End()

	var updatePayload payload.UpdateCategoryRequest
	id, _ := strconv.ParseInt(removeSlashFromParam(c.Param("id")), 10, 64)
	if id == 0 {
		err := fmt.Errorf("id must be a string of numbers")
		c.Error(payload.ErrInvalidRequest(err))
		h.p.Logger.Error("UPDATE_CATEGORY_FAILED", map[string]interface{}{"message": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&updatePayload); err != nil {
		c.Error(payload.ErrInvalidRequest(err))
		h.p.Logger.Error("UPDATE_CATEGORY_FAILED", map[string]interface{}{"message": err.Error()})
		return
	}
	categoryUpdated, err := h.categoryUsecase.UpdateCategoryByID(c, id, updatePayload)
	if err != nil {
		c.Error(err)
		h.p.Logger.Error("UPDATE_CATEGORY_FAILED", map[string]interface{}{"message": err.Error()})
		return
	}
	cateResponse := mapper.CategoryToCategoryResponse(categoryUpdated)
	c.JSON(http.StatusOK, payload.SuccessResponse(cateResponse, ""))
}