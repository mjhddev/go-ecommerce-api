package dto

type CheckoutResponse struct {
	OrderID     uint    `json:"order_id"`
	TotalAmount float64 `json:"total_amount"`
	Status      string  `json:"status"`
}
