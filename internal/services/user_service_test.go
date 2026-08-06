package services_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/mjhddev/go-ecommerce-api/internal/dto"
	"github.com/mjhddev/go-ecommerce-api/internal/errs"
	"github.com/mjhddev/go-ecommerce-api/internal/models"
	"github.com/mjhddev/go-ecommerce-api/internal/services"
	serviceMocks "github.com/mjhddev/go-ecommerce-api/internal/services/mocks"
	"github.com/mjhddev/go-ecommerce-api/internal/utils"
)

func TestRegisterSuccess(t *testing.T) {
	mockRepo := new(serviceMocks.UserRepository)

	service := services.NewUserService(mockRepo)

	req := dto.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	mockRepo.
		On("GetByEmail", req.Email).
		Return(nil, gorm.ErrRecordNotFound)

	mockRepo.
		On("Create", mock.AnythingOfType("*models.User")).
		Run(func(args mock.Arguments) {
			user := args.Get(0).(*models.User)

			user.ID = 1
			user.Role = "user"
		}).
		Return(nil)

	resp, err := service.Register(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	assert.Equal(t, uint(1), resp.ID)
	assert.Equal(t, req.Name, resp.Name)
	assert.Equal(t, req.Email, resp.Email)
	assert.Equal(t, "user", resp.Role)

	mockRepo.AssertExpectations(t)
}

func TestRegisterEmailAlreadyExists(t *testing.T) {
	mockRepo := new(serviceMocks.UserRepository)
	service := services.NewUserService(mockRepo)

	req := dto.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	mockRepo.
		On("GetByEmail", req.Email).
		Return(&models.User{
			ID:    1,
			Name:  "John Doe",
			Email: req.Email,
			Role:  "user",
		}, nil)

	resp, err := service.Register(req)
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrEmailAlreadyExists)

	mockRepo.AssertNotCalled(t, "Create", mock.Anything)

	mockRepo.AssertExpectations(t)
}

func TestRegisterRepositoryError(t *testing.T) {
	mockRepo := new(serviceMocks.UserRepository)

	service := services.NewUserService(mockRepo)

	req := dto.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	dbErr := errors.New("database error")

	mockRepo.
		On("GetByEmail", req.Email).
		Return(nil, dbErr)

	resp, err := service.Register(req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.ErrorIs(t, err, dbErr)

	mockRepo.AssertNotCalled(t, "Create", mock.Anything)
	mockRepo.AssertExpectations(t)
}

func TestRegisterCreateError(t *testing.T) {
	mockRepo := new(serviceMocks.UserRepository)

	service := services.NewUserService(mockRepo)

	req := dto.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	createErr := errors.New("failed to create user")

	mockRepo.
		On("GetByEmail", req.Email).
		Return(nil, gorm.ErrRecordNotFound)

	mockRepo.
		On("Create", mock.AnythingOfType("*models.User")).
		Return(createErr)

	resp, err := service.Register(req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.ErrorIs(t, err, createErr)

	mockRepo.AssertExpectations(t)
}

func TestLoginSuccess(t *testing.T) {
	mockRepo := new(serviceMocks.UserRepository)

	service := services.NewUserService(mockRepo)

	password := "password123"

	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	req := dto.LoginRequest{
		Email:    "john@example.com",
		Password: password,
	}

	mockRepo.
		On("GetByEmail", req.Email).
		Return(&models.User{
			ID:       1,
			Name:     "John",
			Email:    req.Email,
			Password: hashedPassword,
			Role:     "user",
		}, nil)

	resp, err := service.Login(req)

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.NotEmpty(t, resp.AccessToken)

	mockRepo.AssertExpectations(t)
}

func TestLoginUserNotFound(t *testing.T) {
	mockRepo := new(serviceMocks.UserRepository)

	service := services.NewUserService(mockRepo)

	req := dto.LoginRequest{
		Email:    "john@example.com",
		Password: "password123",
	}

	mockRepo.
		On("GetByEmail", req.Email).
		Return(nil, nil)

	resp, err := service.Login(req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrInvalidCredentials)

	mockRepo.AssertExpectations(t)
}

func TestLoginWrongPassword(t *testing.T) {
	mockRepo := new(serviceMocks.UserRepository)

	service := services.NewUserService(mockRepo)

	hashedPassword, _ := utils.HashPassword("correct-password")

	req := dto.LoginRequest{
		Email:    "john@example.com",
		Password: "wrong-password",
	}

	mockRepo.
		On("GetByEmail", req.Email).
		Return(&models.User{
			ID:       1,
			Email:    req.Email,
			Password: hashedPassword,
			Role:     "user",
		}, nil)

	resp, err := service.Login(req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrInvalidCredentials)

	mockRepo.AssertExpectations(t)
}
