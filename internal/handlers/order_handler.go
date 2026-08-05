package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
//	@Success		201	{object}	response.SuccessResponse
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
//	@Success		200	{object}	response.SuccessResponse
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
//	@Success		200	{object}	response.SuccessResponse
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
