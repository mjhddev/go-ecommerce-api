package repositories

import (
	"errors"

	"github.com/mjhddev/go-ecommerce-api/internal/dto"
	"github.com/mjhddev/go-ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Checkout(tx *gorm.DB, order *models.Order) error
	CreateOrderItem(tx *gorm.DB, item *models.OrderItem) error
	GetByUser(userID uint) ([]models.Order, error)
	GetByID(orderID uint) (*models.Order, error)
	GetAll(query dto.AdminOrderQuery) ([]models.Order, error)
	Update(order *models.Order) error
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

func (r *orderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) GetAll(query dto.AdminOrderQuery) ([]models.Order, error) {
	var orders []models.Order

	db := r.db.Model(&models.Order{}).
		Preload("User")

	if query.Search != "" {
		db = db.Joins("JOIN users ON users.id = orders.user_id").
			Where("LOWER(users.name) LIKE LOWER(?)", "%"+query.Search+"%")
	}

	if query.Status != "" {
		db = db.Where("orders.status = ?", query.Status)
	}

	switch query.Sort {
	case "oldest":
		db = db.Order("orders.created_at ASC")
	default:
		db = db.Order("orders.created_at DESC")
	}

	offset := (query.Page - 1) * query.Limit

	err := db.
		Limit(query.Limit).
		Offset(offset).
		Find(&orders).Error

	if err != nil {
		return nil, err
	}

	return orders, nil
}
