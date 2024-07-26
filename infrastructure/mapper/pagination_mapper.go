package mapper

import (
	"pm/domain/entity"
	"pm/infrastructure/controllers/payload"
)

func PaginationToPaginationResponse(pagination *entity.Pagination) payload.PaginationResponse {
	return payload.PaginationResponse{
		Limit:         pagination.Limit,
		Page:          pagination.Page,
		TotalElements: pagination.TotalRows,
		TotalPages:    pagination.TotalPages,
	}
}