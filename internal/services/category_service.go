package services

import (
	"github.com/mjhddev/go-ecommerce-api/internal/dto"
	"github.com/mjhddev/go-ecommerce-api/internal/errs"
	"github.com/mjhddev/go-ecommerce-api/internal/models"
	"github.com/mjhddev/go-ecommerce-api/internal/repositories"
)

type CategoryService interface {
	Create(request dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	GetAll() ([]dto.CategoryResponse, error)
	GetByID(id uint) (*dto.CategoryResponse, error)
}

type categoryService struct {
	repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) CategoryService {
	return &categoryService{
		repo: repo,
	}
}

func (s *categoryService) Create(request dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	category := &models.Category{
		Name: request.Name,
	}

	err := s.repo.Create(category)
	if err != nil {
		return nil, err
	}

	response := &dto.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}

	return response, nil
}

func (s *categoryService) GetAll() ([]dto.CategoryResponse, error) {
	categories, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	response := make([]dto.CategoryResponse, 0, len(categories))
	for _, category := range categories {
		response = append(response, dto.CategoryResponse{
			ID:   category.ID,
			Name: category.Name,
		})
	}

	return response, nil
}

func (s *categoryService) GetByID(id uint) (*dto.CategoryResponse, error) {
	category, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, errs.ErrCategoryNotFound
	}

	response := &dto.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}

	return response, nil
}
