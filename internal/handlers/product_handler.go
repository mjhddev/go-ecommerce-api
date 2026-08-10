package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mjhddev/go-ecommerce-api/internal/dto"
	"github.com/mjhddev/go-ecommerce-api/internal/errs"
	"github.com/mjhddev/go-ecommerce-api/internal/response"
	"github.com/mjhddev/go-ecommerce-api/internal/services"
)

type ProductHandler struct {
	productService services.ProductService
}

func NewProductHandler(productService services.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct godoc
//
//	@Summary		Create product
//	@Description	Create a new product
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.CreateProductRequest	true	"Product Request"
//	@Success		201		{object}	response.SuccessResponse{data=dto.ProductResponse}
//	@Failure		400		{object}	response.ErrorResponse
//	@Failure		401		{object}	response.ErrorResponse
//	@Failure		403		{object}	response.ErrorResponse
//	@Failure		404		{object}	response.ErrorResponse
//	@Router			/products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var request dto.CreateProductRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	product, err := h.productService.Create(request)
	if err != nil {
		if errors.Is(err, errs.ErrCategoryNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"Product created successfully",
		product,
	)
}

// GetAll godoc
//
//	@Summary		Get all products
//	@Description	Get products with pagination
//	@Tags			Products
//	@Produce		json
//	@Security		BearerAuth
//	@Param 			page query 		int 	false "Page Number" default(1)
//	@Param 			limit query 	int 	false "Items per Page" default(10)
//	@Param 			search query 	string 	false "Search product by name"
//	@Param 			category query 	int 	false "Filter by category ID"
//	@Param 			sort query string false "Sort: newest, oldest, price_asc, price_desc, name_asc, name_desc"
//	@Success		200		{object}	response.SuccessResponse{data=dto.ProductListResponse}
//	@Failure		401		{object}	response.ErrorResponse
//	@Failure		500		{object}	response.ErrorResponse
//	@Router			/products [get]
func (h *ProductHandler) GetAll(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	search := strings.TrimSpace(c.Query("search"))

	categoryID64, _ := strconv.ParseUint(c.DefaultQuery("category", "0"), 10, 64)
	categoryID := uint(categoryID64)

	sort := c.DefaultQuery("sort", "newest")

	products, err := h.productService.GetAll(page, limit, search, categoryID, sort)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Products retrieved successfully",
		products,
	)
}

// GetProductByID godoc
//
//	@Summary		Get product by ID
//	@Description	Get product by ID
//	@Tags			Products
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Product ID"
//	@Success		200	{object}	response.SuccessResponse{data=dto.ProductResponse}
//	@Failure		400	{object}	response.ErrorResponse
//	@Failure		404	{object}	response.ErrorResponse
//	@Router			/products/{id} [get]
func (h *ProductHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := h.productService.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, errs.ErrProductNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Product retrieved successfully",
		product,
	)
}

// UpdateProduct godoc
//
//	@Summary		Update product
//	@Description	Update product by ID
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int					true	"Product ID"
//	@Param			request	body		dto.UpdateProductRequest	true	"Product Request"
//	@Success		200		{object}	response.SuccessResponse{data=dto.ProductResponse}
//	@Failure		400		{object}	response.ErrorResponse
//	@Failure		401		{object}	response.ErrorResponse
//	@Failure		403		{object}	response.ErrorResponse
//	@Failure		404		{object}	response.ErrorResponse
//	@Router			/products/{id} [put]
func (h *ProductHandler) Update(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var request dto.UpdateProductRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	product, err := h.productService.Update(uint(id), request)
	if err != nil {
		if errors.Is(err, errs.ErrProductNotFound) || errors.Is(err, errs.ErrCategoryNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Product updated successfully",
		product,
	)
}

// DeleteProduct godoc
//
//	@Summary		Delete product
//	@Description	Delete product by ID
//	@Tags			Products
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Product ID"
//	@Success		200	{object}	response.SuccessResponse
//	@Failure		400	{object}	response.ErrorResponse
//	@Failure		401	{object}	response.ErrorResponse
//	@Failure		403	{object}	response.ErrorResponse
//	@Failure		404	{object}	response.ErrorResponse
//	@Router			/products/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	err = h.productService.Delete(uint(id))
	if err != nil {
		if errors.Is(err, errs.ErrProductNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Product deleted successfully",
		nil,
	)
}
