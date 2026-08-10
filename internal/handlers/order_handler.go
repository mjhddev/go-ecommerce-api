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

type OrderHandler struct {
	orderService services.OrderService
}

func NewOrderHandler(orderService services.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// Checkout godoc
//
//	@Summary		Checkout cart
//	@Description	Create an order from the authenticated user's cart
//	@Tags			Orders
//	@Produce		json
//	@Security		BearerAuth
//	@Success		201	{object}	response.SuccessResponse{data=dto.CheckoutResponse}
//	@Failure		400	{object}	response.ErrorResponse
//	@Failure		401	{object}	response.ErrorResponse
//	@Failure		404	{object}	response.ErrorResponse
//	@Router			/orders/checkout [post]
func (h *OrderHandler) Checkout(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "invalid user id")
		return
	}

	order, err := h.orderService.Checkout(userID)
	if err != nil {

		switch {
		case errors.Is(err, errs.ErrCartEmpty):
			response.Error(c, http.StatusBadRequest, err.Error())

		case errors.Is(err, errs.ErrInsufficientStock):
			response.Error(c, http.StatusBadRequest, err.Error())

		default:
			response.Error(c, http.StatusInternalServerError, err.Error())
		}

		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"Checkout successful",
		order,
	)
}

// GetOrders godoc
//
//	@Summary		Get order history
//	@Description	Get all orders of the authenticated user
//	@Tags			Orders
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.SuccessResponse{data=[]dto.OrderResponse}
//	@Failure		401	{object}	response.ErrorResponse
//	@Router			/orders [get]
func (h *OrderHandler) GetOrders(c *gin.Context) {

	userIDValue, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "invalid user id")
		return
	}

	orders, err := h.orderService.GetOrders(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Orders retrieved successfully",
		orders,
	)
}

// GetOrderByID godoc
//
//	@Summary		Get order detail
//	@Description	Get order detail by ID
//	@Tags			Orders
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Order ID"
//	@Success		200	{object}	response.SuccessResponse{data=dto.OrderDetailResponse}
//	@Failure		400	{object}	response.ErrorResponse
//	@Failure		401	{object}	response.ErrorResponse
//	@Failure		403	{object}	response.ErrorResponse
//	@Failure		404	{object}	response.ErrorResponse
//	@Router			/orders/{id} [get]
func (h *OrderHandler) GetOrderByID(c *gin.Context) {

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid order id")
		return
	}

	userID := c.MustGet("userID").(uint)

	order, err := h.orderService.GetOrderByID(userID, uint(id))
	if err != nil {

		switch {

		case errors.Is(err, errs.ErrOrderNotFound):
			response.Error(c, http.StatusNotFound, err.Error())

		case errors.Is(err, errs.ErrForbidden):
			response.Error(c, http.StatusForbidden, err.Error())

		default:
			response.Error(c, http.StatusInternalServerError, err.Error())
		}

		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Order retrieved successfully",
		order,
	)
}

// GetAllAdminOrders godoc
//
//	@Summary		Get all orders
//	@Description	Get all orders for admin
//	@Tags			Admin Orders
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page	query	int	false	"Page"
//	@Param			limit	query	int	false	"Limit"
//	@Param			search	query	string	false	"Customer name"
//	@Param			status	query	string	false	"Order status"
//	@Param			sort	query	string	false	"newest or oldest"
//	@Success 		200 {object} response.SuccessResponse{data=[]dto.AdminOrderResponse}
//	@Router			/admin/orders [get]
func (h *OrderHandler) GetAllAdminOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	query := dto.AdminOrderQuery{
		Page:   page,
		Limit:  limit,
		Search: strings.TrimSpace(c.Query("search")),
		Status: strings.TrimSpace(c.Query("status")),
		Sort:   c.DefaultQuery("sort", "newest"),
	}

	orders, err := h.orderService.GetAllAdminOrders(query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Orders retrieved successfully",
		orders,
	)
}

// GetAdminOrderByID godoc
//
//	@Summary		Get order detail
//	@Description	Get order detail by id
//	@Tags			Admin Orders
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Order ID"
//	@Success 		200 {object} 	response.SuccessResponse{data=dto.AdminOrderDetailResponse}
//	@Failure		404	{object}	response.ErrorResponse
//	@Router			/admin/orders/{id} [get]
func (h *OrderHandler) GetAdminOrderByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	order, err := h.orderService.GetAdminOrderByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Order retrieved successfully",
		order,
	)
}

// UpdateOrderStatus godoc
//
//	@Summary		Update order status
//	@Description	Update order status by admin
//	@Tags			Admin Orders
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path	int								true	"Order ID"
//	@Param			request	body	dto.UpdateOrderStatusRequest	true	"Request Body"
//	@Success		200		{object}	response.SuccessResponse{data=dto.AdminOrderDetailResponse}
//	@Failure		400		{object}	response.ErrorResponse
//	@Failure		401		{object}	response.ErrorResponse
//	@Failure		403		{object}	response.ErrorResponse
//	@Failure		404		{object}	response.ErrorResponse
//	@Failure		500		{object}	response.ErrorResponse
//	@Router			/admin/orders/{id}/status [put]
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var request dto.UpdateOrderStatusRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.orderService.UpdateOrderStatus(uint(id), request)
	if err != nil {
		switch err {
		case errs.ErrOrderNotFound:
			response.Error(c, http.StatusNotFound, err.Error())
		case errs.ErrInvalidOrderStatus:
			response.Error(c, http.StatusBadRequest, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Order status updated successfully",
		order,
	)
}
