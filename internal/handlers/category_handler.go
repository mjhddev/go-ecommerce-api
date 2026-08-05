package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mjhddev/go-ecommerce-api/internal/dto"
	"github.com/mjhddev/go-ecommerce-api/internal/errs"
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

// CreateCategory godoc
//
//	@Summary		Create category
//	@Description	Create a new category
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.CreateCategoryRequest	true	"Category Request"
//	@Success		201		{object}	response.SuccessResponse
//	@Failure		400		{object}	response.ErrorResponse
//	@Failure		401		{object}	response.ErrorResponse
//	@Failure		403		{object}	response.ErrorResponse
//	@Router			/categories [post]
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

// GetCategories godoc
//
//	@Summary		Get all categories
//	@Description	Get all categories
//	@Tags			Categories
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.SuccessResponse
//	@Failure		401	{object}	response.ErrorResponse
//	@Router			/categories [get]
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

// GetCategoryByID godoc
//
//	@Summary		Get category by ID
//	@Description	Get category by ID
//	@Tags			Categories
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Category ID"
//	@Success		200	{object}	response.SuccessResponse
//	@Failure		400	{object}	response.ErrorResponse
//	@Failure		404	{object}	response.ErrorResponse
//	@Router			/categories/{id} [get]
func (h *CategoryHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	category, err := h.categoryService.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, errs.ErrCategoryNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
		} else {
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Category retrieved successfully",
		category,
	)
}

// UpdateCategory godoc
//
//	@Summary		Update category
//	@Description	Update category by ID
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int					true	"Category ID"
//	@Param			request	body		dto.UpdateCategoryRequest	true	"Category Request"
//	@Success		200		{object}	response.SuccessResponse
//	@Failure		400		{object}	response.ErrorResponse
//	@Failure		401		{object}	response.ErrorResponse
//	@Failure		403		{object}	response.ErrorResponse
//	@Failure		404		{object}	response.ErrorResponse
//	@Router			/categories/{id} [put]
func (h *CategoryHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var request dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.categoryService.Update(uint(id), request)
	if err != nil {
		if errors.Is(err, errs.ErrCategoryNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
		} else {
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Category updated successfully",
		category,
	)
}

// DeleteCategory godoc
//
//	@Summary		Delete category
//	@Description	Delete category by ID
//	@Tags			Categories
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Category ID"
//	@Success		200	{object}	response.SuccessResponse
//	@Failure		400	{object}	response.ErrorResponse
//	@Failure		401	{object}	response.ErrorResponse
//	@Failure		403	{object}	response.ErrorResponse
//	@Failure		404	{object}	response.ErrorResponse
//	@Router			/categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	err = h.categoryService.Delete(uint(id))
	if err != nil {
		if errors.Is(err, errs.ErrCategoryNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
		} else {
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Category deleted successfully",
		nil,
	)
}
