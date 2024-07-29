package mapper

import (
	"pm/domain/entity"
	"pm/infrastructure/controllers/payload"
)

func CategoryToCategoryResponse(e *entity.Category) payload.CategoryResponse {
	return payload.CategoryResponse{
		ID:   e.ID,
		Name: e.Name,
		AuditTime: payload.AuditTime{
			UpdatedAt: e.UpdatedAt,
			CreatedAt: e.CreatedAt,
		},
	}
}

func CategoriesToListCategoriesResponse(listEntities []entity.Category, pagination *entity.Pagination) payload.ListCategoryResponses {
	listCateResponse := make([]payload.CategoryResponse, 0)
	for _, c := range listEntities {
		cateResponse := CategoryToCategoryResponse(&c)
		listCateResponse = append(listCateResponse, cateResponse)
	}

	return payload.ListCategoryResponses{
		Categories:         listCateResponse,
		PaginationResponse: PaginationToPaginationResponse(pagination),
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