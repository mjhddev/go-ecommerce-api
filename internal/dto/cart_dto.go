package dto

type AddToCartRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required"`
}

type CartItemResponse struct {
	ID       uint            `json:"id"`
	Product  ProductResponse `json:"product"`
	Quantity int             `json:"quantity"`
}

type UpdateCartRequest struct {
	Quantity int `json:"quantity" binding:"required,gte=1"`
}
