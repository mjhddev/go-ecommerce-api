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

type CartHandler struct {
	cartService services.CartService
}

func NewCartHandler(cartService services.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

func (h *CartHandler) AddToCart(c *gin.Context) {
	var request dto.AddToCartRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "invalid user ID")
		return
	}

	cartItem, err := h.cartService.AddToCart(userID, request)
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
		http.StatusCreated, "product added to cart successfully",
		cartItem,
	)

}

func (h *CartHandler) GetCart(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "invalid user ID")
		return
	}

	cartItems, err := h.cartService.GetCart(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(
		c,
		http.StatusOK, "cart items retrieved successfully",
		cartItems,
	)
}

func (h *CartHandler) Update(c *gin.Context) {
	id := c.Param("id")

	cartID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid cart item ID")
		return
	}

	var request dto.UpdateCartRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userIDValue, _ := c.Get("userID")
	userID := userIDValue.(uint)

	cartItem, err := h.cartService.Update(userID, uint(cartID), request)
	if err != nil {

		switch {
		case errors.Is(err, errs.ErrCartItemNotFound):
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
		http.StatusOK, "cart item updated successfully",
		cartItem,
	)
}

func (h *CartHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	cartID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid cart item ID")
		return
	}

	userIDValue, _ := c.Get("userID")
	userID := userIDValue.(uint)

	err = h.cartService.Delete(userID, uint(cartID))
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrCartItemNotFound):
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
		http.StatusOK, "cart item deleted successfully",
		nil,
	)
}
