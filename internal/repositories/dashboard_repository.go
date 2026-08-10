package repositories

import (
	"github.com/mjhddev/go-ecommerce-api/internal/dto"
	"github.com/mjhddev/go-ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type DashboardRepository interface {
	GetTotalOrders() (int64, error)
	GetTotalRevenue() (float64, error)
	GetPendingOrders() (int64, error)
	GetCompletedOrders() (int64, error)

	GetMonthlyRevenue() ([]dto.MonthlyRevenueResponse, error)

	GetTopProducts(limit int) ([]dto.TopProductResponse, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{
		db: db,
	}
}

func (r *dashboardRepository) GetTotalOrders() (int64, error) {
	var total int64

	err := r.db.Model(&models.Order{}).
		Count(&total).Error

	return total, err
}

func (r *dashboardRepository) GetTotalRevenue() (float64, error) {
	var revenue float64

	err := r.db.Model(&models.Order{}).
		Where("status = ?", "completed").
		Select("COALESCE(SUM(total_amount),0)").
		Scan(&revenue).Error

	return revenue, err
}

func (r *dashboardRepository) GetPendingOrders() (int64, error) {
	var total int64

	err := r.db.Model(&models.Order{}).
		Where("status = ?", "pending").
		Count(&total).Error

	return total, err
}

func (r *dashboardRepository) GetCompletedOrders() (int64, error) {
	var total int64

	err := r.db.Model(&models.Order{}).
		Where("status = ?", "completed").
		Count(&total).Error

	return total, err
}

func (r *dashboardRepository) GetMonthlyRevenue() ([]dto.MonthlyRevenueResponse, error) {

	var revenue []dto.MonthlyRevenueResponse

	err := r.db.
		Model(&models.Order{}).
		Select(`
			TO_CHAR(created_at, 'FMMonth') AS month,
			COALESCE(SUM(total_amount),0) AS revenue
		`).
		Where("status = ?", "completed").
		Group("TO_CHAR(created_at, 'FMMonth'), EXTRACT(MONTH FROM created_at)").
		Order("EXTRACT(MONTH FROM created_at)").
		Scan(&revenue).Error

	return revenue, err
}

func (r *dashboardRepository) GetTopProducts(limit int) ([]dto.TopProductResponse, error) {

	var products []dto.TopProductResponse

	err := r.db.
		Table("order_items").
		Select(`
			products.id,
			products.name,
			products.image_url,
			COALESCE(SUM(order_items.quantity),0) AS sold
		`).
		Joins("JOIN products ON products.id = order_items.product_id").
		Group("products.id").
		Order("sold DESC").
		Limit(limit).
		Scan(&products).Error

	return products, err
}
