package controllers

import "time"

// ------------------ General Error Response ------------------ //

type ErrorResponse struct {
	Error string `json:"error"`
}

// ------------------ Auth Response ------------------ //

type UserPayload struct {
	ID      uint   `json:"id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
}

// RegisterResponse is returned when a user registers successfully
type RegisterResponse struct {
	Message string      `json:"message"`
	User    UserPayload `json:"user"`
}

// LoginResponse is returned when user logs in successfully
type LoginResponse struct {
	Token string `json:"token"`
}

// ------------------ Product Response ------------------ //

type ProductPayload struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateProductResponse is returned after creating a product
type CreateProductResponse struct {
	Data ProductPayload `json:"data"`
}

// GetProductsResponse is returned when listing all products
type GetProductsResponse struct {
	Data []ProductPayload `json:"data"`
}

// SingleProductResponse is returned when getting a single product
type SingleProductResponse struct {
	Data ProductPayload `json:"data"`
}

// UpdateProductResponse is returned after updating a product
type UpdateProductResponse struct {
	Data ProductPayload `json:"data"`
}

// DeleteProductResponse is a simple message for deletion success
type DeleteProductResponse struct {
	Message string `json:"message"`
}

// ------------------ Order Response ------------------ //

type OrderItemPayload struct {
	ProductID   uint    `json:"product_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

type OrderPayload struct {
	ID        uint               `json:"id"`
	UserID    uint               `json:"user_id"`
	Status    string             `json:"status"`
	Products  []OrderItemPayload `json:"products"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

// CreateOrderResponse is returned after creating an order
type CreateOrderResponse struct {
	Data OrderPayload `json:"data"`
}

// GetOrdersResponse is returned when listing user’s orders
type GetOrdersResponse struct {
	Data []OrderPayload `json:"data"`
}

// CancelOrderResponse is returned after cancelling an order
type CancelOrderResponse struct {
	Message string `json:"message"`
}

// UpdateOrderStatusResponse is returned after updating an order’s status
type UpdateOrderStatusResponse struct {
	Data OrderPayload `json:"data"`
}
