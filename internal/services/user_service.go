package services

import (
	"errors"

	"github.com/mjhddev/go-ecommerce-api/internal/dto"
	"github.com/mjhddev/go-ecommerce-api/internal/errs"
	"github.com/mjhddev/go-ecommerce-api/internal/models"
	"github.com/mjhddev/go-ecommerce-api/internal/repositories"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Register(request dto.RegisterRequest) (*dto.RegisterResponse, error)
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

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: string(hashedPassword),
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
