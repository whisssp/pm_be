package payload

type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gte=0"`
	CategoryID  int64   `json:"categoryId" validate:"required"`
	Stock       int64   `json:"stock" validate:"gte=0"`
	Image       string  `json:"imagePath"`
}

type UpdateProductRequest struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gte=0"`
	CategoryID  int64   `json:"categoryId" validate:"required"`
	Stock       int64   `json:"stock" validate:"gte=0"`
	Image       string  `json:"imagePath"`
}