package payload

import "time"

type AuditTime struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

type PaginationResponse struct {
	Limit         int   `json:"limit"`
	Page          int   `json:"page"`
	TotalElements int64 `json:"totalElements"`
	TotalPages    int   `json:"totalPages"`
}

type ProductResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryID  int64   `json:"categoryId"`
	Stock       int64   `json:"stock"`
	Image       string  `json:"imagePath"`
	AuditTime
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
	AuditTime
}

type ListCategoryResponses struct {
	Categories []CategoryResponse `json:"categories"`
	PaginationResponse
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Role     string `json:"role"`
	AuditTime
}

type ListUserResponses struct {
	Users []UserResponse `json:"users"`
	PaginationResponse
}

type AuthResponse struct {
	Token string `json:"token"`
}

type OrderResponse struct {
	ID         int64               `json:"id"`
	UserID     int64               `json:"userId"`
	Status     string              `json:"status"`
	OrderItems []OrderItemResponse `json:"orderItems"`
	Total      float64             `json:"total"`
	AuditTime
}

type OrderItemResponse struct {
	ID        int64   `json:"id"`
	ProductID int64   `json:"productId"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	AuditTime
}

type ListOrderResponses struct {
	Orders []OrderResponse `json:"orders"`
	PaginationResponse
}