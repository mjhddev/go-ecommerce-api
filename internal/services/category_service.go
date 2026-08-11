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
	Update(id uint, request dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
	Delete(id uint) error
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

func (s *categoryService) Update(id uint, request dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	category, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, errs.ErrCategoryNotFound
	}

	category.Name = request.Name

	err = s.repo.Update(category)
	if err != nil {
		return nil, err
	}

	response := &dto.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}

	return response, nil
}

func (s *categoryService) Delete(id uint) error {
	category, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if category == nil {
		return errs.ErrCategoryNotFound
	}

	return s.repo.Delete(id)
}
