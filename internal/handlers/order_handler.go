package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mjhddev/go-ecommerce-api/internal/errs"
	"github.com/mjhddev/go-ecommerce-api/internal/response"
	"github.com/mjhddev/go-ecommerce-api/internal/services"
)

type OrderHanlder struct {
	orderService services.OrderService
}

func NewOrderHandler(orderService services.OrderService) *OrderHanlder {
	return &OrderHanlder{
		orderService: orderService,
	}
}

func (h *OrderHanlder) Checkout(c *gin.Context) {
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
