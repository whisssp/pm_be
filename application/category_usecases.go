package application

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"pm/domain/entity"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/implementations/categories"
	"pm/infrastructure/mapper"
	"pm/infrastructure/persistences/base"
	"pm/utils"
)

type CategoryUsecase interface {
	CreateCategory(*gin.Context, *payload.CreateCategoryRequest) error
	GetAllCategories(*gin.Context, *entity.CategoryFilter, *entity.Pagination) (*payload.ListCategoryResponses, error)
	GetCategoryByID(*gin.Context, int64) (*payload.CategoryResponse, error)
	DeleteCategoryByID(*gin.Context, int64) error
	UpdateCategoryByID(*gin.Context, int64, payload.UpdateCategoryRequest) (*payload.CategoryResponse, error)
}

type categoryUsecase struct {
	p *base.Persistence
}

func NewCategoryUsecase(p *base.Persistence) CategoryUsecase {
	return categoryUsecase{p}
}

func (categoryUsecase categoryUsecase) UpdateCategoryByID(c *gin.Context, id int64, updatePayload payload.UpdateCategoryRequest) (*payload.CategoryResponse, error) {
	span := categoryUsecase.p.Logger.Start(c, "UPDATE_CATEGORY_USECASES", categoryUsecase.p.Logger.SetContextWithSpanFunc())
	defer span.End()
	categoryUsecase.p.Logger.Info("UPDATE_CATEGORY", map[string]interface{}{"data": updatePayload})

	categoryRepo := categories.NewCategoryRepository(c, categoryUsecase.p, categoryUsecase.p.GormDB)
	cate, err := categoryRepo.GetCategoryByID(id)
	if err != nil {
		categoryUsecase.p.Logger.Error("UPDATE_CATEGORY_FAILED", map[string]interface{}{"message": err.Error()})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, payload.ErrEntityNotFound("categories", err)
		}
		return nil, payload.ErrInvalidRequest(err)
	}
	mapper.UpdateCategory(cate, &updatePayload)
	cate, err = categoryRepo.Update(cate)
	if err != nil {
		categoryUsecase.p.Logger.Error("UPDATE_CATEGORY_FAILED", map[string]interface{}{"message": err.Error()})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, payload.ErrEntityNotFound("categories", err)
		}
		return nil, payload.ErrInvalidRequest(err)
	}
	cateResponse := mapper.CategoryToCategoryResponse(cate)
	categoryUsecase.p.Logger.Info("UPDATE_CATEGORY_SUCCESSFULLY", map[string]interface{}{"data": cateResponse})
	return &cateResponse, nil
}

func (categoryUsecase categoryUsecase) GetCategoryByID(c *gin.Context, id int64) (*payload.CategoryResponse, error) {
	span := categoryUsecase.p.Logger.Start(c, "GET_CATEGORY_USECASES", categoryUsecase.p.Logger.SetContextWithSpanFunc())
	defer span.End()
	categoryUsecase.p.Logger.Info("GET_CATEGORY", map[string]interface{}{"data": id})

	categoryRepo := categories.NewCategoryRepository(c, categoryUsecase.p, categoryUsecase.p.GormDB)
	cate, err := categoryRepo.GetCategoryByID(id)
	if err != nil {
		categoryUsecase.p.Logger.Error("GET_CATEGORY_FAILED", map[string]interface{}{"message": err.Error()})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, payload.ErrEntityNotFound("categories", err)
		}
		return nil, payload.ErrInvalidRequest(err)
	}
	cateResponse := mapper.CategoryToCategoryResponse(cate)

	categoryUsecase.p.Logger.Info("GET_CATEGORY_SUCCESSFULLY", map[string]interface{}{"data": cateResponse})
	return &cateResponse, nil
}

func (categoryUsecase categoryUsecase) DeleteCategoryByID(c *gin.Context, id int64) error {
	span := categoryUsecase.p.Logger.Start(c, "DELETE_CATEGORY_USECASES", categoryUsecase.p.Logger.SetContextWithSpanFunc())
	defer span.End()
	categoryUsecase.p.Logger.Info("DELETE_CATEGORY", map[string]interface{}{"data": id})

	categoryRepo := categories.NewCategoryRepository(c, categoryUsecase.p, categoryUsecase.p.GormDB)
	cate, err := categoryRepo.GetCategoryByID(id)
	if err != nil {
		categoryUsecase.p.Logger.Error("DELETE_CATEGORY_FAILED", map[string]interface{}{"message": err.Error()})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return payload.ErrEntityNotFound("categories", err)
		}
		return payload.ErrInvalidRequest(err)
	}
	if err := categoryRepo.DeleteCategory(cate); err != nil {
		categoryUsecase.p.Logger.Error("DELETE_CATEGORY_FAILED", map[string]interface{}{"message": err.Error()})
		return payload.ErrDB(err)
	}

	categoryUsecase.p.Logger.Info("CREATE_CATEGORY_SUCCESSFULLY", map[string]interface{}{"data": cate.ID})
	return nil
}

func (categoryUsecase categoryUsecase) CreateCategory(c *gin.Context, reqPayload *payload.CreateCategoryRequest) error {
	span := categoryUsecase.p.Logger.Start(c, "CREATE_CATEGORY_USECASES", categoryUsecase.p.Logger.SetContextWithSpanFunc())
	defer span.End()
	categoryUsecase.p.Logger.Info("CREATE_CATEGORY", map[string]interface{}{"data": reqPayload})

	if err := utils.ValidateReqPayload(reqPayload); err != nil {
		categoryUsecase.p.Logger.Error("CREATE_CATEGORY_FAILED", map[string]interface{}{"message": err.Error()})
		return payload.ErrValidateFailed(err)
	}

	categoryEntity := mapper.CreateCatePayloadToCategory(reqPayload)
	cateRepo := categories.NewCategoryRepository(c, categoryUsecase.p, categoryUsecase.p.GormDB)
	err := cateRepo.Create(categoryEntity)
	if err != nil {
		categoryUsecase.p.Logger.Error("CREATE_CATEGORY_FAILED", map[string]interface{}{"message": err.Error()})
		return err
	}

	categoryUsecase.p.Logger.Info("CREATE_CATEGORY_SUCCESSFULLY", map[string]interface{}{"data": categoryEntity.ID})
	return nil
}

func (categoryUsecase categoryUsecase) GetAllCategories(c *gin.Context, filter *entity.CategoryFilter, pagination *entity.Pagination) (*payload.ListCategoryResponses, error) {
	span := categoryUsecase.p.Logger.Start(c, "GET_ALL_CATEGORIES_USECASES", categoryUsecase.p.Logger.SetContextWithSpanFunc())
	defer span.End()
	categoryUsecase.p.Logger.Info("GET_ALL_CATEGORIES", map[string]interface{}{
		"params": struct {
			Filter     *entity.CategoryFilter `json:"filter"`
			Pagination *entity.Pagination     `json:"pagination"`
		}{
			Filter:     filter,
			Pagination: pagination,
		},
	})

	cateRepo := categories.NewCategoryRepository(c, categoryUsecase.p, categoryUsecase.p.GormDB)
	cates, err := cateRepo.GetAllCategories(filter, pagination)
	if err != nil {
		categoryUsecase.p.Logger.Info("GET_ALL_CATEGORIES_FAILED", map[string]interface{}{"message": err.Error()})
		return nil, err
	}
	listCatesResponse := mapper.CategoriesToListCategoriesResponse(cates, pagination)
	categoryUsecase.p.Logger.Info("GET_ALL_CATEGORIES_SUCCESSFULLY", map[string]interface{}{"data": listCatesResponse})
	return &listCatesResponse, nil
}