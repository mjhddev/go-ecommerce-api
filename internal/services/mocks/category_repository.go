package mocks

import (
	"github.com/mjhddev/go-ecommerce-api/internal/models"
	"github.com/stretchr/testify/mock"
)

type CategoryRepository struct {
	mock.Mock
}

func (m *CategoryRepository) Create(category *models.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *CategoryRepository) GetAll() ([]models.Category, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]models.Category), args.Error(1)
}

func (m *CategoryRepository) GetByID(id uint) (*models.Category, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *CategoryRepository) Update(category *models.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *CategoryRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
