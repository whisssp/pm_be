package payload

import "time"

type ProductResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CategoryID  int64     `json:"categoryId"`
	Stock       int64     `json:"stock"`
	Image       string    `json:"imagePath"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
}

type ListProductResponse struct {
	Products      []ProductResponse `json:"products"`
	Limit         int               `json:"limit"`
	Page          int               `json:"page"`
	TotalElements int64             `json:"totalElements"`
	TotalPages    int               `json:"totalPages"`
}

type CategoryResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

type ListCategoriesResponse struct {
	Categories    []CategoryResponse `json:"categories"`
	Limit         int                `json:"limit"`
	Page          int                `json:"page"`
	TotalElements int64              `json:"totalElements"`
	TotalPages    int                `json:"totalPages"`
}