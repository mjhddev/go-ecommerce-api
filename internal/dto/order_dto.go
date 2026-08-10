package dto

type CheckoutResponse struct {
	OrderID     uint    `json:"order_id"`
	TotalAmount float64 `json:"total_amount"`
	Status      string  `json:"status"`
}

type OrderResponse struct {
	ID          uint    `json:"id"`
	TotalAmount float64 `json:"total_amount"`
	Status      string  `json:"status"`
}

type OrderItemResponse struct {
	Product  OrderProductResponse `json:"product"`
	Quantity int                  `json:"quantity"`
	Price    float64              `json:"price"`
}

type OrderProductResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type OrderDetailResponse struct {
	ID          uint                `json:"id"`
	TotalAmount float64             `json:"total_amount"`
	Status      string              `json:"status"`
	Items       []OrderItemResponse `json:"items"`
}
