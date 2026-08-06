package services_test

import (
	"errors"
	"testing"

	"github.com/mjhddev/go-ecommerce-api/internal/errs"
	"github.com/mjhddev/go-ecommerce-api/internal/models"
	"github.com/mjhddev/go-ecommerce-api/internal/services"
	serviceMocks "github.com/mjhddev/go-ecommerce-api/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetOrdersSuccess(t *testing.T) {
	orderRepo := new(serviceMocks.OrderRepository)
	cartRepo := new(serviceMocks.CartRepository)
	productRepo := new(serviceMocks.ProductRepository)

	service := services.NewOrderService(
		nil,
		orderRepo,
		cartRepo,
		productRepo,
	)

	orders := []models.Order{
		{
			ID:          1,
			UserID:      1,
			TotalAmount: 1000000,
			Status:      "pending",
		},
		{
			ID:          2,
			UserID:      1,
			TotalAmount: 500000,
			Status:      "paid",
		},
	}

	orderRepo.
		On("GetByUser", uint(1)).
		Return(orders, nil)

	resp, err := service.GetOrders(1)

	require.NoError(t, err)
	require.Len(t, resp, 2)

	assert.Equal(t, uint(1), resp[0].ID)
	assert.Equal(t, 1000000.0, resp[0].TotalAmount)
	assert.Equal(t, "pending", resp[0].Status)

	assert.Equal(t, uint(2), resp[1].ID)
	assert.Equal(t, "paid", resp[1].Status)

	orderRepo.AssertExpectations(t)
}

func TestGetOrdersRepositoryError(t *testing.T) {
	orderRepo := new(serviceMocks.OrderRepository)
	cartRepo := new(serviceMocks.CartRepository)
	productRepo := new(serviceMocks.ProductRepository)

	service := services.NewOrderService(
		nil,
		orderRepo,
		cartRepo,
		productRepo,
	)

	dbErr := errors.New("database error")

	orderRepo.
		On("GetByUser", uint(1)).
		Return(nil, dbErr)

	resp, err := service.GetOrders(1)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, dbErr)

	orderRepo.AssertExpectations(t)
}

func TestGetOrderByIDSuccess(t *testing.T) {
	orderRepo := new(serviceMocks.OrderRepository)
	cartRepo := new(serviceMocks.CartRepository)
	productRepo := new(serviceMocks.ProductRepository)

	service := services.NewOrderService(
		nil,
		orderRepo,
		cartRepo,
		productRepo,
	)

	order := &models.Order{
		ID:          1,
		UserID:      1,
		TotalAmount: 1100000,
		Status:      "pending",
		OrderItems: []models.OrderItem{
			{
				ProductID: 1,
				Quantity:  2,
				Price:     750000,
				Product: models.Product{
					ID:   1,
					Name: "Keyboard",
				},
			},
			{
				ProductID: 2,
				Quantity:  1,
				Price:     350000,
				Product: models.Product{
					ID:   2,
					Name: "Mouse",
				},
			},
		},
	}

	orderRepo.
		On("GetByID", uint(1)).
		Return(order, nil)

	resp, err := service.GetOrderByID(1, 1)

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, uint(1), resp.ID)
	assert.Equal(t, "pending", resp.Status)
	assert.Len(t, resp.Items, 2)

	assert.Equal(t, "Keyboard", resp.Items[0].Product.Name)
	assert.Equal(t, int(2), resp.Items[0].Quantity)

	orderRepo.AssertExpectations(t)
}

func TestGetOrderByIDNotFound(t *testing.T) {
	orderRepo := new(serviceMocks.OrderRepository)
	cartRepo := new(serviceMocks.CartRepository)
	productRepo := new(serviceMocks.ProductRepository)

	service := services.NewOrderService(
		nil,
		orderRepo,
		cartRepo,
		productRepo,
	)

	orderRepo.
		On("GetByID", uint(1)).
		Return(nil, nil)

	resp, err := service.GetOrderByID(1, 1)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, errs.ErrOrderNotFound)

	orderRepo.AssertExpectations(t)
}

func TestGetOrderByIDForbidden(t *testing.T) {
	orderRepo := new(serviceMocks.OrderRepository)
	cartRepo := new(serviceMocks.CartRepository)
	productRepo := new(serviceMocks.ProductRepository)

	service := services.NewOrderService(
		nil,
		orderRepo,
		cartRepo,
		productRepo,
	)

	orderRepo.
		On("GetByID", uint(1)).
		Return(&models.Order{
			ID:     1,
			UserID: 99,
		}, nil)

	resp, err := service.GetOrderByID(1, 1)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, errs.ErrForbidden)

	orderRepo.AssertExpectations(t)
}

func TestGetOrderByIDRepositoryError(t *testing.T) {
	orderRepo := new(serviceMocks.OrderRepository)
	cartRepo := new(serviceMocks.CartRepository)
	productRepo := new(serviceMocks.ProductRepository)

	service := services.NewOrderService(
		nil,
		orderRepo,
		cartRepo,
		productRepo,
	)

	dbErr := errors.New("database error")

	orderRepo.
		On("GetByID", uint(1)).
		Return(nil, dbErr)

	resp, err := service.GetOrderByID(1, 1)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, dbErr)

	orderRepo.AssertExpectations(t)
}
