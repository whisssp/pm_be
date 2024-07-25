package mapper

import (
	"pm/domain/entity"
	"pm/infrastructure/controllers/payload"
)

func CategoryToCategoryResponse(e *entity.Category) payload.CategoryResponse {
	return payload.CategoryResponse{
		ID:        e.ID,
		Name:      e.Name,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func CategoriesToListCategoriesResponse(listEntities []entity.Category, pagination *entity.Pagination) payload.ListCategoriesResponse {
	listCateResponse := make([]payload.CategoryResponse, 0)
	for _, c := range listEntities {
		cateResponse := CategoryToCategoryResponse(&c)
		listCateResponse = append(listCateResponse, cateResponse)
	}
	return payload.ListCategoriesResponse{
		Categories:    listCateResponse,
		Limit:         pagination.Limit,
		Page:          pagination.Page,
		TotalElements: pagination.TotalRows,
		TotalPages:    pagination.TotalPages,
	}
}

func CreateCatePayloadToCategory(reqPayload *payload.CreateCategoryRequest) *entity.Category {
	return &entity.Category{
		Name: reqPayload.Name,
	}
}

func UpdateCategory(old *entity.Category, updatePayload *payload.UpdateCategoryRequest) {
	old.Name = updatePayload.Name

}