package repositories

import (
	"errors"

	"github.com/mjhddev/go-ecommerce-api/internal/dto"
	"github.com/mjhddev/go-ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *models.Product) error
	GetAll(query dto.ProductQuery) ([]models.Product, int64, error)
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

func (r *productRepository) GetAll(query dto.ProductQuery) ([]models.Product, int64, error) {
	var (
		products []models.Product
		total    int64
	)

	offset := (query.Page - 1) * query.Limit

	db := r.db.Model(&models.Product{})

	if query.Search != "" {
		db = db.Where("LOWER(name) LIKE LOWER(?)", "%"+query.Search+"%")
	}

	if query.CategoryID != 0 {
		db = db.Where("category_id = ?", query.CategoryID)
	}

	switch query.Sort {
	case "price_asc":
		db = db.Order("price ASC")
	case "price_desc":
		db = db.Order("price DESC")
	case "name_asc":
		db = db.Order("name ASC")
	case "name_desc":
		db = db.Order("name DESC")
	case "newest":
		db = db.Order("created_at DESC")
	case "oldest":
		db = db.Order("created_at ASC")
	default:
		db = db.Order("id DESC")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.
		Preload("Category").
		Limit(query.Limit).
		Offset(offset).
		Find(&products).Error

	if err != nil {
		return nil, 0, err
	}
	return products, total, err
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
