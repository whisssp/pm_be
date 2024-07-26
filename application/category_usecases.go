package application

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"pm/domain/entity"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/implementations/categories"
	"pm/infrastructure/mapper"
	"pm/infrastructure/persistences/base"
	"pm/utils"
)

type CategoryUsecase interface {
	CreateCategory(reqPayload *payload.CreateCategoryRequest) error
	GetAllCategories(filter *entity.CategoryFilter, pagination *entity.Pagination) (*payload.ListCategoryResponses, error)
	GetCategoryByID(id int64) (*payload.CategoryResponse, error)
	DeleteCategoryByID(id int64) error
	UpdateCategoryByID(id int64, updatePayload payload.UpdateCategoryRequest) (*payload.CategoryResponse, error)
}

type categoryUsecase struct {
	p *base.Persistence
}

func NewCategoryUsecase(p *base.Persistence) CategoryUsecase {
	return categoryUsecase{p}
}

func (categoryUsecase categoryUsecase) UpdateCategoryByID(id int64, updatePayload payload.UpdateCategoryRequest) (*payload.CategoryResponse, error) {
	categoryRepo := categories.NewCategoryRepository(categoryUsecase.p)
	cate, err := categoryRepo.GetCategoryByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, payload.ErrEntityNotFound("categories", err)
		}
		return nil, payload.ErrInvalidRequest(err)
	}
	mapper.UpdateCategory(cate, &updatePayload)
	cateResponse := mapper.CategoryToCategoryResponse(cate)
	return &cateResponse, nil
}

func (categoryUsecase categoryUsecase) GetCategoryByID(id int64) (*payload.CategoryResponse, error) {
	categoryRepo := categories.NewCategoryRepository(categoryUsecase.p)
	cate, err := categoryRepo.GetCategoryByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, payload.ErrEntityNotFound("categories", err)
		}
		return nil, payload.ErrInvalidRequest(err)
	}
	cateResponse := mapper.CategoryToCategoryResponse(cate)
	return &cateResponse, nil
}

func (categoryUsecase categoryUsecase) DeleteCategoryByID(id int64) error {
	categoryRepo := categories.NewCategoryRepository(categoryUsecase.p)
	cate, err := categoryRepo.GetCategoryByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return payload.ErrEntityNotFound("categories", err)
		}
		return payload.ErrInvalidRequest(err)
	}
	if err := categoryRepo.DeleteCategory(cate); err != nil {
		return payload.ErrInvalidRequest(err)
	}
	return nil
}

func (categoryUsecase categoryUsecase) CreateCategory(reqPayload *payload.CreateCategoryRequest) error {
	if err := utils.ValidateReqPayload(reqPayload); err != nil {
		fmt.Printf("error validating categories: %v", err)
		return payload.ErrValidateFailed(err)
	}

	categoryEntity := mapper.CreateCatePayloadToCategory(reqPayload)
	cateRepo := categories.NewCategoryRepository(categoryUsecase.p)
	err := cateRepo.Create(categoryEntity)
	if err != nil {
		fmt.Printf("error creating categories: %v", err)
		return payload.ErrDB(err)
	}

	return nil
}

func (categoryUsecase categoryUsecase) GetAllCategories(filter *entity.CategoryFilter, pagination *entity.Pagination) (*payload.ListCategoryResponses, error) {
	cateRepo := categories.NewCategoryRepository(categoryUsecase.p)
	cates, err := cateRepo.GetAllCategories(filter, pagination)
	if err != nil {
		return nil, payload.ErrDB(err)
	}
	listCatesResponse := mapper.CategoriesToListCategoriesResponse(cates, pagination)
	return &listCatesResponse, nil
}