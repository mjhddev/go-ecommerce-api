package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mjhddev/go-ecommerce-api/internal/dto"
	"github.com/mjhddev/go-ecommerce-api/internal/errs"
	"github.com/mjhddev/go-ecommerce-api/internal/response"
	"github.com/mjhddev/go-ecommerce-api/internal/services"
)

type AuthHandler struct {
	userService services.UserService
}

func NewAuthHandler(userService services.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var request dto.RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userService.Register(request)
	if err != nil {
		if errors.Is(err, errs.ErrEmailAlreadyExists) {
			response.Error(c, http.StatusConflict, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"User registered successfully",
		user,
	)
}
