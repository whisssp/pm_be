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

	categoryRepo := categories.NewCategoryRepository(c, categoryUsecase.p, categoryUsecase.p.GormDB)
	cate, err := categoryRepo.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}
	mapper.UpdateCategory(cate, &updatePayload)
	cate, err = categoryRepo.Update(cate)
	if err != nil {
		return nil, err
	}

	categoryUsecase.p.Logger.Info("UPDATE_CATEGORY_SUCCESSFULLY", map[string]interface{}{"category_response": cate})
	return cate, nil
}

func (categoryUsecase categoryUsecase) GetCategoryByID(c *gin.Context, id int64) (*entity.Category, error) {

	categoryRepo := categories.NewCategoryRepository(c, categoryUsecase.p, categoryUsecase.p.GormDB)
	cate, err := categoryRepo.GetCategoryByID(id)
	if err != nil {
		categoryUsecase.p.Logger.Error("GET_CATEGORY_FAILED", map[string]interface{}{"message": err.Error()})
		return nil, err
	}

	categoryUsecase.p.Logger.Info("GET_CATEGORY_SUCCESSFULLY", map[string]interface{}{"data": cate})
	return cate, nil
}

func (categoryUsecase categoryUsecase) DeleteCategoryByID(c *gin.Context, id int64) error {
	categoryRepo := categories.NewCategoryRepository(c, categoryUsecase.p, categoryUsecase.p.GormDB)
	cate, err := categoryRepo.GetCategoryByID(id)
	if err != nil {
		return err
	}
	if err := categoryRepo.DeleteCategory(cate); err != nil {
		return err
	}

	return nil
}

func (categoryUsecase categoryUsecase) CreateCategory(c *gin.Context, reqPayload *payload.CreateCategoryRequest) error {

	if err := utils.ValidateReqPayload(reqPayload); err != nil {
		return payload.ErrValidateFailed(err)
	}

	categoryEntity := mapper.CreateCatePayloadToCategory(reqPayload)
	cateRepo := categories.NewCategoryRepository(c, categoryUsecase.p, categoryUsecase.p.GormDB)
	err := cateRepo.Create(categoryEntity)
	if err != nil {
		return err
	}

	return nil
}

func (categoryUsecase categoryUsecase) GetAllCategories(c *gin.Context, filter *entity.CategoryFilter, pagination *entity.Pagination) ([]entity.Category, error) {
	cateRepo := categories.NewCategoryRepository(c, categoryUsecase.p, categoryUsecase.p.GormDB)
	cates, err := cateRepo.GetAllCategories(filter, pagination)
	if err != nil {
		return nil, err
	}

	return cates, nil
}