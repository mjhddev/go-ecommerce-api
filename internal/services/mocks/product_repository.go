package mocks

import (
	"github.com/mjhddev/go-ecommerce-api/internal/models"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type ProductRepository struct {
	mock.Mock
}

func (m *ProductRepository) Create(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *ProductRepository) GetAll() ([]models.Product, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *ProductRepository) GetByID(id uint) (*models.Product, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *ProductRepository) Update(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *ProductRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *ProductRepository) UpdateTx(tx *gorm.DB, product *models.Product) error {
	args := m.Called(tx, product)
	return args.Error(0)
}
