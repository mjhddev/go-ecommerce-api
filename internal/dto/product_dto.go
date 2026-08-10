package dto

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"required,gte=0"`
	CategoryID  uint    `json:"category_id" binding:"required"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"required,gte=0"`
	CategoryID  uint    `json:"category_id" binding:"required"`
}

type ProductResponse struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Price       float64          `json:"price"`
	Stock       int              `json:"stock"`
	ImageURL    string           `json:"image_url"`
	Category    CategoryResponse `json:"category"`
}

type PaginationResponse struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	TotalData int64 `json:"total_data"`
	TotalPage int   `json:"total_page"`
}

type ProductListResponse struct {
	Items      []ProductResponse  `json:"items"`
	Pagination PaginationResponse `json:"pagination"`
}
