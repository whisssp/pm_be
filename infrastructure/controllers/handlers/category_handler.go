package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"pm/application"
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

func (h CategoryHandler) HandleGetAllCategories(c *gin.Context) {
	var categoryFilter payload.CategoryFilter
	var pagination payload.Pagination

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