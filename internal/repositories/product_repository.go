package repositories

import (
	"errors"

	"github.com/mjhddev/go-ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *models.Product) error
	GetAll() ([]models.Product, error)
	GetByID(id uint) (*models.Product, error)
	Update(product *models.Product) error
	Delete(id uint) error
	UpdateTx(tx *gorm.DB, product *models.Product) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) GetAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Category").Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, err
}

func (r *productRepository) GetByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Category").First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

func (r *productRepository) UpdateTx(tx *gorm.DB, product *models.Product) error {
	return tx.Save(product).Error
}
