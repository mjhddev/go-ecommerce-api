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

func (h *AuthHandler) Login(c *gin.Context) {
	var request dto.LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	loginResponse, err := h.userService.Login(request)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidCredentials) {
			response.Error(c, http.StatusUnauthorized, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}
	response.Success(c, http.StatusOK, "Login successful", loginResponse)
}

func (h *AuthHandler) Profile(c *gin.Context) {
	userID := c.GetUint("userID")
	profile, err := h.userService.Profile(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Profile retrieved successfully",
		profile,
	)
}
