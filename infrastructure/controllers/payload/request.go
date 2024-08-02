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

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateCategoryRequest struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type UserRequest struct {
	Name     string `json:"name" validate:"required,max=150"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,max=11,e164"`
	Password string `json:"password" validate:"required,min=6,max=11"`
	Role     int64  `json:"role" validate:"oneof=1 2,required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=11"`
}

type UpdateOrderRequest struct {
}

type CreateOrderRequest struct {
	OrderItems []OrderItemRequest `json:"orderItems"`
	UserID     uint               `json:"userId"`
	Status     string             `json:"status"`
}

type OrderItemRequest struct {
	ProductID uint    `json:"productId"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}