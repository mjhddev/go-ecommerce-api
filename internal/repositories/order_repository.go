package repositories

import (
	"github.com/mjhddev/go-ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Checkout(tx *gorm.DB, order *models.Order) error
	CreateOrderItem(tx *gorm.DB, item *models.OrderItem) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Checkout(tx *gorm.DB, order *models.Order) error {
	return tx.Create(order).Error
}

func (r *orderRepository) CreateOrderItem(tx *gorm.DB, item *models.OrderItem) error {
	return tx.Create(item).Error
}
