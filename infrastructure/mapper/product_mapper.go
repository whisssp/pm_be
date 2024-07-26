package mapper

import (
	"pm/domain/entity"
	"pm/infrastructure/controllers/payload"
)

func ProductToProductResponse(product *entity.Product) payload.ProductResponse {
	return payload.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CategoryID:  product.CategoryID,
		Stock:       product.Stock,
		Image:       product.Image,
		UpdatedAt:   product.UpdatedAt,
		CreatedAt:   product.CreatedAt,
	}
}

func ProdsToListProdsResponse(products []entity.Product, pagination *entity.Pagination) payload.ListProductResponses {
	listProdResponse := make([]payload.ProductResponse, 0)
	for _, p := range products {
		prodResponse := ProductToProductResponse(&p)
		listProdResponse = append(listProdResponse, prodResponse)
	}
	return payload.ListProductResponses{
		Products:           listProdResponse,
		PaginationResponse: PaginationToPaginationResponse(pagination),
	}
}

func PayloadToProduct(reqPayload *payload.CreateProductRequest) *entity.Product {
	return &entity.Product{
		Name:        reqPayload.Name,
		Description: reqPayload.Description,
		Price:       reqPayload.Price,
		CategoryID:  reqPayload.CategoryID,
		Stock:       reqPayload.Stock,
		Image:       reqPayload.Image,
	}
}

func UpdateProduct(oldProd *entity.Product, updatePayload *payload.UpdateProductRequest) {
	oldProd.Name = updatePayload.Name
	oldProd.Description = updatePayload.Description
	oldProd.Stock = updatePayload.Stock
	oldProd.Price = updatePayload.Price
	oldProd.Image = updatePayload.Image
}