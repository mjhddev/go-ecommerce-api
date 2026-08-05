package repositories

import (
	"errors"

	"github.com/mjhddev/go-ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type CartRepository interface {
	Create(cart *models.CartItem) error
	GetByUserAndProduct(userID, productID uint) (*models.CartItem, error)
	Update(cart *models.CartItem) error
	GetByUser(userID uint) ([]models.CartItem, error)
	GetByID(id uint) (*models.CartItem, error)
	Delete(id uint) error
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) Create(cart *models.CartItem) error {
	return r.db.Create(cart).Error
}

func (r *cartRepository) GetByUserAndProduct(userID, productID uint) (*models.CartItem, error) {
	var cart models.CartItem
	err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) Update(cart *models.CartItem) error {
	return r.db.Save(cart).Error
}

func (r *cartRepository) GetByUser(userID uint) ([]models.CartItem, error) {
	var cartItems []models.CartItem
	err := r.db.
		Preload("Product").
		Preload("Product.Category").
		Where("user_id = ?", userID).
		Find(&cartItems).Error

	if err != nil {
		return nil, err
	}
	return cartItems, nil
}

func (r *cartRepository) GetByID(id uint) (*models.CartItem, error) {
	var cartItem models.CartItem
	err := r.db.Preload("Product").Preload("Product.Category").First(&cartItem, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cartItem, nil
}

func (r *cartRepository) Delete(id uint) error {
	return r.db.Delete(&models.CartItem{}, id).Error
}
