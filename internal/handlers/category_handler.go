package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mjhddev/go-ecommerce-api/internal/dto"
	"github.com/mjhddev/go-ecommerce-api/internal/response"
	"github.com/mjhddev/go-ecommerce-api/internal/services"
)

type CategoryHandler struct {
	categoryService services.CategoryService
}

func NewCategoryHandler(categoryService services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var request dto.CreateCategoryRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.categoryService.Create(request)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"Category created successfully",
		category,
	)
}

func (h *CategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.categoryService.GetAll()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Categories retrieved successfully",
		categories,
	)
}
