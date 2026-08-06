package mocks

import (
	"github.com/mjhddev/go-ecommerce-api/internal/models"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type OrderRepository struct {
	mock.Mock
}

func (m *OrderRepository) Checkout(tx *gorm.DB, order *models.Order) error {
	args := m.Called(tx, order)
	return args.Error(0)
}

func (m *OrderRepository) CreateOrderItem(tx *gorm.DB, item *models.OrderItem) error {
	args := m.Called(tx, item)
	return args.Error(0)
}

func (m *OrderRepository) GetByUser(userID uint) ([]models.Order, error) {
	args := m.Called(userID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]models.Order), args.Error(1)
}

func (m *OrderRepository) GetByID(orderID uint) (*models.Order, error) {
	args := m.Called(orderID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.Order), args.Error(1)
}
