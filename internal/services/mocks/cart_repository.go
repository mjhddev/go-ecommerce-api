package mocks

import (
	"github.com/mjhddev/go-ecommerce-api/internal/models"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type CartRepository struct {
	mock.Mock
}

func (m *CartRepository) Create(cart *models.CartItem) error {
	args := m.Called(cart)
	return args.Error(0)
}

func (m *CartRepository) GetByUserAndProduct(userID, productID uint) (*models.CartItem, error) {
	args := m.Called(userID, productID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.CartItem), args.Error(1)
}

func (m *CartRepository) Update(cart *models.CartItem) error {
	args := m.Called(cart)
	return args.Error(0)
}

func (m *CartRepository) GetByUser(userID uint) ([]models.CartItem, error) {
	args := m.Called(userID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]models.CartItem), args.Error(1)
}

func (m *CartRepository) GetByID(id uint) (*models.CartItem, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.CartItem), args.Error(1)
}

func (m *CartRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *CartRepository) DeleteByUser(tx *gorm.DB, userID uint) error {
	args := m.Called(tx, userID)
	return args.Error(0)
}
