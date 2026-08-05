package repositories

import (
	"errors"

	"github.com/mjhddev/go-ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Checkout(tx *gorm.DB, order *models.Order) error
	CreateOrderItem(tx *gorm.DB, item *models.OrderItem) error
	GetByUser(userID uint) ([]models.Order, error)
	GetByID(orderID uint) (*models.Order, error)
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

func (r *orderRepository) GetByUser(userID uint) ([]models.Order, error) {
	var orders []models.Order

	err := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *orderRepository) GetByID(orderID uint) (*models.Order, error) {
	var order models.Order

	err := r.db.
		Preload("OrderItems").
		Preload("OrderItems.Product").
		First(&order, orderID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &order, nil
}
