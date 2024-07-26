package payload

import "time"

type PaginationResponse struct {
	Limit         int   `json:"limit"`
	Page          int   `json:"page"`
	TotalElements int64 `json:"totalElements"`
	TotalPages    int   `json:"totalPages"`
}

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

type ListProductResponses struct {
	Products []ProductResponse `json:"products"`
	PaginationResponse
}

type CategoryResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

type ListCategoryResponses struct {
	Categories []CategoryResponse `json:"categories"`
	PaginationResponse
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

type ListUserResponses struct {
	Users []UserResponse `json:"users"`
	PaginationResponse
}

type AuthResponse struct {
	Token string `json:"token"`
}