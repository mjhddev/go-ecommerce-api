package services

import (
	"github.com/mjhddev/go-ecommerce-api/internal/dto"
	"github.com/mjhddev/go-ecommerce-api/internal/repositories"
)

type DashboardService interface {
	GetDashboard() (*dto.DashboardResponse, error)
}

type dashboardService struct {
	dashboardRepo repositories.DashboardRepository
}

func NewDashboardService(
	dashboardRepo repositories.DashboardRepository,
) DashboardService {
	return &dashboardService{
		dashboardRepo: dashboardRepo,
	}
}

func (s *dashboardService) GetDashboard() (*dto.DashboardResponse, error) {

	totalOrders, err := s.dashboardRepo.GetTotalOrders()
	if err != nil {
		return nil, err
	}

	totalRevenue, err := s.dashboardRepo.GetTotalRevenue()
	if err != nil {
		return nil, err
	}

	pendingOrders, err := s.dashboardRepo.GetPendingOrders()
	if err != nil {
		return nil, err
	}

	completedOrders, err := s.dashboardRepo.GetCompletedOrders()
	if err != nil {
		return nil, err
	}

	monthlyRevenue, err := s.dashboardRepo.GetMonthlyRevenue()
	if err != nil {
		return nil, err
	}

	topProducts, err := s.dashboardRepo.GetTopProducts(5)
	if err != nil {
		return nil, err
	}

	response := &dto.DashboardResponse{
		TotalOrders:     totalOrders,
		TotalRevenue:    totalRevenue,
		PendingOrders:   pendingOrders,
		CompletedOrders: completedOrders,
		MonthlyRevenue:  monthlyRevenue,
		TopProducts:     topProducts,
	}

	return response, nil
}
