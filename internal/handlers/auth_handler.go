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

// Register godoc
//
//	@Summary		Register a new user
//	@Description	Register a new user account
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RegisterRequest	true	"Register Request"
//	@Success		201		{object}	response.SuccessResponse{data=dto.RegisterResponse}
//	@Failure		400		{object}	response.ErrorResponse
//	@Failure		409		{object}	response.ErrorResponse
//	@Router			/auth/register [post]
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

// Login godoc
//
//	@Summary		Login user
//	@Description	Authenticate user and return JWT token
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.LoginRequest	true	"Login Request"
//	@Success		200		{object}	response.SuccessResponse{data=dto.LoginResponse}
//	@Failure		400		{object}	response.ErrorResponse
//	@Failure		401		{object}	response.ErrorResponse
//	@Router			/auth/login [post]
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

// Profile godoc
//
//	@Summary		Get user profile
//	@Description	Get authenticated user profile
//	@Tags			User
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.SuccessResponse{data=dto.ProfileResponse}
//	@Failure		401	{object}	response.ErrorResponse
//	@Router			/profile [get]
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
