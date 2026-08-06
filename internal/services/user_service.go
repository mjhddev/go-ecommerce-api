package services

import (
	"errors"

	"github.com/mjhddev/go-ecommerce-api/internal/dto"
	"github.com/mjhddev/go-ecommerce-api/internal/errs"
	"github.com/mjhddev/go-ecommerce-api/internal/models"
	"github.com/mjhddev/go-ecommerce-api/internal/repositories"
	"github.com/mjhddev/go-ecommerce-api/internal/utils"
	"gorm.io/gorm"
)

type UserService interface {
	Register(request dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(request dto.LoginRequest) (*dto.LoginResponse, error)
	Profile(userID uint) (*dto.ProfileResponse, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) Register(request dto.RegisterRequest) (*dto.RegisterResponse, error) {
	existingUser, err := s.repo.GetByEmail(request.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, errs.ErrEmailAlreadyExists
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: hashedPassword,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	response := &dto.RegisterResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	return response, nil
}

func (s *userService) Login(request dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repo.GetByEmail(request.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errs.ErrInvalidCredentials
	}

	if err := utils.CheckPassword(user.Password, request.Password); err != nil {
		return nil, errs.ErrInvalidCredentials
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken: token,
	}, nil
}

func (s *userService) Profile(userID uint) (*dto.ProfileResponse, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errs.ErrUserNotFound
	}

	return &dto.ProfileResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}
