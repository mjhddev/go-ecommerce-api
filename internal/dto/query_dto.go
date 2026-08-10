package dto

type ProductQuery struct {
	Page       int    `form:"page"`
	Limit      int    `form:"limit"`
	Search     string `form:"search"`
	CategoryID uint   `form:"category"`
	Sort       string `form:"sort"`
}
