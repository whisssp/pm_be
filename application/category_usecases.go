package application

import (
	"github.com/gin-gonic/gin"
	"pm/domain/entity"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/implementations/categories"
	"pm/infrastructure/mapper"
	"pm/infrastructure/persistences/base"
	"pm/utils"
)

type CategoryUsecase interface {
	CreateCategory(*gin.Context, *payload.CreateCategoryRequest) error
	GetAllCategories(*gin.Context, *entity.CategoryFilter, *entity.Pagination) ([]entity.Category, error)
	GetCategoryByID(*gin.Context, int64) (*entity.Category, error)
	DeleteCategoryByID(*gin.Context, int64) error
	UpdateCategoryByID(*gin.Context, int64, payload.UpdateCategoryRequest) (*entity.Category, error)
}

type categoryUsecase struct {
	p *base.Persistence
}

func NewCategoryUsecase(p *base.Persistence) CategoryUsecase {
	return categoryUsecase{p}
}

func (categoryUsecase categoryUsecase) UpdateCategoryByID(c *gin.Context, id int64, updatePayload payload.UpdateCategoryRequest) (*entity.Category, error) {
	span := categoryUsecase.p.Logger.Start(c, "UPDATE_CATEGORY_USECASES", categoryUsecase.p.Logger.SetContextWithSpanFunc())
	defer span.End()
	categoryUsecase.p.Logger.Info("UPDATE_CATEGORY", map[string]interface{}{"data": updatePayload})

	categoryRepo := categories.NewCategoryRepository(c, categoryUsecase.p, categoryUsecase.p.GormDB)
	cate, err := categoryRepo.GetCategoryByID(span, id)
	if err != nil {
		categoryUsecase.p.Logger.Error("UPDATE_CATEGORY: ERROR", map[string]interface{}{"error": err.Error()})
		return nil, err
	}
	mapper.UpdateCategory(cate, &updatePayload)
	cate, err = categoryRepo.Update(span, cate)
	if err != nil {
		categoryUsecase.p.Logger.Error("UPDATE_CATEGORY: ERROR", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	categoryUsecase.p.Logger.Info("UPDATE_CATEGORY_SUCCESSFULLY", map[string]interface{}{"category_response": cate})
	return cate, nil
}

func (categoryUsecase categoryUsecase) GetCategoryByID(c *gin.Context, id int64) (*entity.Category, error) {
	span := categoryUsecase.p.Logger.Start(c, "GET_CATEGORY_USECASES", categoryUsecase.p.Logger.SetContextWithSpanFunc())
	defer span.End()
	categoryUsecase.p.Logger.Info("GET_CATEGORY", map[string]interface{}{"data": id})

	categoryRepo := categories.NewCategoryRepository(c, categoryUsecase.p, categoryUsecase.p.GormDB)
	cate, err := categoryRepo.GetCategoryByID(span, id)
	if err != nil {
		categoryUsecase.p.Logger.Error("GET_CATEGORY_FAILED", map[string]interface{}{"message": err.Error()})
		return nil, err
	}

	categoryUsecase.p.Logger.Info("GET_CATEGORY_SUCCESSFULLY", map[string]interface{}{"data": cate})
	return cate, nil
}

func (categoryUsecase categoryUsecase) DeleteCategoryByID(c *gin.Context, id int64) error {
	span := categoryUsecase.p.Logger.Start(c, "DELETE_CATEGORY_USECASES", categoryUsecase.p.Logger.SetContextWithSpanFunc())
	defer span.End()
	categoryUsecase.p.Logger.Info("DELETE_CATEGORY", map[string]interface{}{"data": id})

	categoryRepo := categories.NewCategoryRepository(c, categoryUsecase.p, categoryUsecase.p.GormDB)
	cate, err := categoryRepo.GetCategoryByID(span, id)
	if err != nil {
		categoryUsecase.p.Logger.Error("DELETE_CATEGORY_FAILED", map[string]interface{}{"message": err.Error()})
		return err
	}
	if err := categoryRepo.DeleteCategory(span, cate); err != nil {
		categoryUsecase.p.Logger.Error("DELETE_CATEGORY_FAILED", map[string]interface{}{"message": err.Error()})
		return err
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
	err := cateRepo.Create(span, categoryEntity)
	if err != nil {
		categoryUsecase.p.Logger.Error("CREATE_CATEGORY_FAILED", map[string]interface{}{"message": err.Error()})
		return err
	}

	categoryUsecase.p.Logger.Info("CREATE_CATEGORY_SUCCESSFULLY", map[string]interface{}{"data": categoryEntity.ID})
	return nil
}

func (categoryUsecase categoryUsecase) GetAllCategories(c *gin.Context, filter *entity.CategoryFilter, pagination *entity.Pagination) ([]entity.Category, error) {
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
	cates, err := cateRepo.GetAllCategories(span, filter, pagination)
	if err != nil {
		categoryUsecase.p.Logger.Info("GET_ALL_CATEGORIES_FAILED", map[string]interface{}{"message": err.Error()})
		return nil, err
	}

	categoryUsecase.p.Logger.Info("GET_ALL_CATEGORIES_SUCCESSFULLY", map[string]interface{}{"data": cates})
	return cates, nil
}