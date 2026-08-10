package dto

import "time"

type AdminOrderQuery struct {
	Page   int
	Limit  int
	Search string
	Status string
	Sort   string
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type AdminOrderResponse struct {
	ID          uint      `json:"id"`
	CustomerID  uint      `json:"customer_id"`
	Customer    string    `json:"customer"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type AdminOrderDetailResponse struct {
	ID          uint                `json:"id"`
	CustomerID  uint                `json:"customer_id"`
	Customer    string              `json:"customer"`
	TotalAmount float64             `json:"total_amount"`
	Status      string              `json:"status"`
	CreatedAt   time.Time           `json:"created_at"`
	Items       []OrderItemResponse `json:"items"`
}
