package services_test

import (
	"errors"
	"testing"

	"github.com/mjhddev/go-ecommerce-api/internal/dto"
	"github.com/mjhddev/go-ecommerce-api/internal/errs"
	"github.com/mjhddev/go-ecommerce-api/internal/models"
	"github.com/mjhddev/go-ecommerce-api/internal/services"
	serviceMocks "github.com/mjhddev/go-ecommerce-api/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAddToCartSuccess(t *testing.T) {
	cartRepo := new(serviceMocks.CartRepository)
	productRepo := new(serviceMocks.ProductRepository)

	service := services.NewCartService(cartRepo, productRepo)

	req := dto.AddToCartRequest{
		ProductID: 1,
		Quantity:  2,
	}

	product := &models.Product{
		ID:    1,
		Name:  "Keyboard",
		Price: 750000,
		Stock: 10,
		Category: models.Category{
			ID:   2,
			Name: "Gaming",
		},
	}

	productRepo.
		On("GetByID", uint(1)).
		Return(product, nil)

	cartRepo.
		On("GetByUserAndProduct", uint(1), uint(1)).
		Return(nil, nil)

	cartRepo.
		On("Create", mock.AnythingOfType("*models.CartItem")).
		Run(func(args mock.Arguments) {
			cart := args.Get(0).(*models.CartItem)
			cart.ID = 100
		}).
		Return(nil)

	resp, err := service.AddToCart(1, req)

	require.NoError(t, err)

	assert.Equal(t, uint(100), resp.ID)
	assert.Equal(t, int(2), resp.Quantity)
	assert.Equal(t, "Keyboard", resp.Product.Name)

	cartRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}

func TestAddToCartUpdateExisting(t *testing.T) {
	cartRepo := new(serviceMocks.CartRepository)
	productRepo := new(serviceMocks.ProductRepository)

	service := services.NewCartService(cartRepo, productRepo)

	req := dto.AddToCartRequest{
		ProductID: 1,
		Quantity:  3,
	}

	product := &models.Product{
		ID:   1,
		Name: "Keyboard",
		Category: models.Category{
			ID:   2,
			Name: "Gaming",
		},
	}

	cart := &models.CartItem{
		ID:        10,
		UserID:    1,
		ProductID: 1,
		Quantity:  2,
	}

	productRepo.
		On("GetByID", uint(1)).
		Return(product, nil)

	cartRepo.
		On("GetByUserAndProduct", uint(1), uint(1)).
		Return(cart, nil)

	cartRepo.
		On("Update", mock.AnythingOfType("*models.CartItem")).
		Run(func(args mock.Arguments) {
			item := args.Get(0).(*models.CartItem)
			assert.Equal(t, uint(1), item.UserID)
			assert.Equal(t, uint(1), item.ProductID)
			assert.Equal(t, int(5), item.Quantity)
		}).
		Return(nil)

	resp, err := service.AddToCart(1, req)

	require.NoError(t, err)

	assert.Equal(t, int(5), resp.Quantity)

	cartRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}

func TestAddToCartProductNotFound(t *testing.T) {
	cartRepo := new(serviceMocks.CartRepository)
	productRepo := new(serviceMocks.ProductRepository)

	service := services.NewCartService(cartRepo, productRepo)

	req := dto.AddToCartRequest{
		ProductID: 1,
		Quantity:  2,
	}

	productRepo.
		On("GetByID", uint(1)).
		Return(nil, nil)

	resp, err := service.AddToCart(1, req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, errs.ErrProductNotFound)

	cartRepo.AssertNotCalled(t, "GetByUserAndProduct", mock.Anything, mock.Anything)
	cartRepo.AssertNotCalled(t, "Create", mock.Anything)
	cartRepo.AssertNotCalled(t, "Update", mock.Anything)

	productRepo.AssertExpectations(t)
	cartRepo.AssertExpectations(t)
}

func TestAddToCartProductRepositoryError(t *testing.T) {
	cartRepo := new(serviceMocks.CartRepository)
	productRepo := new(serviceMocks.ProductRepository)

	service := services.NewCartService(cartRepo, productRepo)

	req := dto.AddToCartRequest{
		ProductID: 1,
		Quantity:  2,
	}

	dbErr := errors.New("database error")

	productRepo.
		On("GetByID", uint(1)).
		Return(nil, dbErr)

	resp, err := service.AddToCart(1, req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, dbErr)

	cartRepo.AssertNotCalled(t, "GetByUserAndProduct", mock.Anything, mock.Anything)

	productRepo.AssertExpectations(t)
	cartRepo.AssertExpectations(t)
}

func TestAddToCartGetCartError(t *testing.T) {
	cartRepo := new(serviceMocks.CartRepository)
	productRepo := new(serviceMocks.ProductRepository)

	service := services.NewCartService(cartRepo, productRepo)

	req := dto.AddToCartRequest{
		ProductID: 1,
		Quantity:  2,
	}

	dbErr := errors.New("database error")

	product := &models.Product{
		ID:   1,
		Name: "Keyboard",
	}

	productRepo.
		On("GetByID", uint(1)).
		Return(product, nil)

	cartRepo.
		On("GetByUserAndProduct", uint(1), uint(1)).
		Return(nil, dbErr)

	resp, err := service.AddToCart(1, req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, dbErr)

	cartRepo.AssertNotCalled(t, "Create", mock.Anything)
	cartRepo.AssertNotCalled(t, "Update", mock.Anything)

	productRepo.AssertExpectations(t)
	cartRepo.AssertExpectations(t)
}

func TestAddToCartCreateError(t *testing.T) {
	cartRepo := new(serviceMocks.CartRepository)
	productRepo := new(serviceMocks.ProductRepository)

	service := services.NewCartService(cartRepo, productRepo)

	req := dto.AddToCartRequest{
		ProductID: 1,
		Quantity:  2,
	}

	createErr := errors.New("create failed")

	product := &models.Product{
		ID:   1,
		Name: "Keyboard",
		Category: models.Category{
			ID:   2,
			Name: "Gaming",
		},
	}

	productRepo.
		On("GetByID", uint(1)).
		Return(product, nil)

	cartRepo.
		On("GetByUserAndProduct", uint(1), uint(1)).
		Return(nil, nil)

	cartRepo.
		On("Create", mock.AnythingOfType("*models.CartItem")).
		Return(createErr)

	resp, err := service.AddToCart(1, req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, createErr)

	productRepo.AssertExpectations(t)
	cartRepo.AssertExpectations(t)
}

func TestAddToCartUpdateError(t *testing.T) {
	cartRepo := new(serviceMocks.CartRepository)
	productRepo := new(serviceMocks.ProductRepository)

	service := services.NewCartService(cartRepo, productRepo)

	req := dto.AddToCartRequest{
		ProductID: 1,
		Quantity:  2,
	}

	updateErr := errors.New("update failed")

	product := &models.Product{
		ID:   1,
		Name: "Keyboard",
		Category: models.Category{
			ID:   2,
			Name: "Gaming",
		},
	}

	cart := &models.CartItem{
		ID:        1,
		UserID:    1,
		ProductID: 1,
		Quantity:  5,
	}

	productRepo.
		On("GetByID", uint(1)).
		Return(product, nil)

	cartRepo.
		On("GetByUserAndProduct", uint(1), uint(1)).
		Return(cart, nil)

	cartRepo.
		On("Update", mock.AnythingOfType("*models.CartItem")).
		Return(updateErr)

	resp, err := service.AddToCart(1, req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, updateErr)

	productRepo.AssertExpectations(t)
	cartRepo.AssertExpectations(t)
}

func TestGetCartSuccess(t *testing.T) {
	cartRepo := new(serviceMocks.CartRepository)
	productRepo := new(serviceMocks.ProductRepository)

	service := services.NewCartService(cartRepo, productRepo)

	cartItems := []models.CartItem{
		{
			ID:       1,
			UserID:   1,
			Quantity: 2,
			Product: models.Product{
				ID:          10,
				Name:        "Mechanical Keyboard",
				Description: "RGB Keyboard",
				Price:       750000,
				Stock:       20,
				Category: models.Category{
					ID:   2,
					Name: "Gaming",
				},
			},
		},
		{
			ID:       2,
			UserID:   1,
			Quantity: 1,
			Product: models.Product{
				ID:          11,
				Name:        "Gaming Mouse",
				Description: "Wireless Mouse",
				Price:       350000,
				Stock:       15,
				Category: models.Category{
					ID:   2,
					Name: "Gaming",
				},
			},
		},
	}

	cartRepo.
		On("GetByUser", uint(1)).
		Return(cartItems, nil)

	resp, err := service.GetCart(1)

	require.NoError(t, err)
	require.Len(t, resp, 2)

	assert.Equal(t, uint(1), resp[0].ID)
	assert.Equal(t, int(2), resp[0].Quantity)
	assert.Equal(t, "Mechanical Keyboard", resp[0].Product.Name)

	assert.Equal(t, uint(2), resp[1].ID)
	assert.Equal(t, "Gaming Mouse", resp[1].Product.Name)

	cartRepo.AssertExpectations(t)
}

func TestGetCartEmpty(t *testing.T) {
	cartRepo := new(serviceMocks.CartRepository)
	productRepo := new(serviceMocks.ProductRepository)

	service := services.NewCartService(cartRepo, productRepo)

	cartRepo.
		On("GetByUser", uint(1)).
		Return([]models.CartItem{}, nil)

	resp, err := service.GetCart(1)

	require.NoError(t, err)

	assert.Empty(t, resp)

	cartRepo.AssertExpectations(t)
}

func TestGetCartRepositoryError(t *testing.T) {
	cartRepo := new(serviceMocks.CartRepository)
	productRepo := new(serviceMocks.ProductRepository)

	service := services.NewCartService(cartRepo, productRepo)

	dbErr := errors.New("database error")

	cartRepo.
		On("GetByUser", uint(1)).
		Return(nil, dbErr)

	resp, err := service.GetCart(1)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, dbErr)

	cartRepo.AssertExpectations(t)
}

func TestUpdateCartSuccess(t *testing.T) {
	cartRepo := new(serviceMocks.CartRepository)
	productRepo := new(serviceMocks.ProductRepository)

	service := services.NewCartService(cartRepo, productRepo)

	req := dto.UpdateCartRequest{
		Quantity: 5,
	}

	cart := &models.CartItem{
		ID:       1,
		UserID:   1,
		Quantity: 2,
		Product: models.Product{
			ID:          10,
			Name:        "Mechanical Keyboard",
			Description: "RGB Keyboard",
			Price:       750000,
			Stock:       20,
			Category: models.Category{
				ID:   2,
				Name: "Gaming",
			},
		},
	}

	cartRepo.
		On("GetByID", uint(1)).
		Return(cart, nil)

	cartRepo.
		On("Update", mock.AnythingOfType("*models.CartItem")).
		Run(func(args mock.Arguments) {
			item := args.Get(0).(*models.CartItem)
			assert.Equal(t, int(5), item.Quantity)
		}).
		Return(nil)

	resp, err := service.Update(1, 1, req)

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, int(5), resp.Quantity)
	assert.Equal(t, "Mechanical Keyboard", resp.Product.Name)

	cartRepo.AssertExpectations(t)
}

func TestUpdateCartNotFound(t *testing.T) {
	cartRepo := new(serviceMocks.CartRepository)
	productRepo := new(serviceMocks.ProductRepository)

	service := services.NewCartService(cartRepo, productRepo)

	req := dto.UpdateCartRequest{
		Quantity: 5,
	}

	cartRepo.
		On("GetByID", uint(1)).
		Return(nil, nil)

	resp, err := service.Update(1, 1, req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, errs.ErrCartItemNotFound)

	cartRepo.AssertNotCalled(t, "Update", mock.Anything)

	cartRepo.AssertExpectations(t)
}

func TestUpdateCartForbidden(t *testing.T) {
	cartRepo := new(serviceMocks.CartRepository)
	productRepo := new(serviceMocks.ProductRepository)

	service := services.NewCartService(cartRepo, productRepo)

	req := dto.UpdateCartRequest{
		Quantity: 5,
	}

	cartRepo.
		On("GetByID", uint(1)).
		Return(&models.CartItem{
			ID:     1,
			UserID: 99,
		}, nil)

	resp, err := service.Update(1, 1, req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, errs.ErrForbidden)

	cartRepo.AssertNotCalled(t, "Update", mock.Anything)

	cartRepo.AssertExpectations(t)
}

func TestUpdateCartRepositoryError(t *testing.T) {
	cartRepo := new(serviceMocks.CartRepository)
	productRepo := new(serviceMocks.ProductRepository)

	service := services.NewCartService(cartRepo, productRepo)

	req := dto.UpdateCartRequest{
		Quantity: 5,
	}

	dbErr := errors.New("database error")

	cartRepo.
		On("GetByID", uint(1)).
		Return(nil, dbErr)

	resp, err := service.Update(1, 1, req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, dbErr)

	cartRepo.AssertNotCalled(t, "Update", mock.Anything)

	cartRepo.AssertExpectations(t)
}

func TestUpdateCartUpdateError(t *testing.T) {
	cartRepo := new(serviceMocks.CartRepository)
	productRepo := new(serviceMocks.ProductRepository)

	service := services.NewCartService(cartRepo, productRepo)

	req := dto.UpdateCartRequest{
		Quantity: 5,
	}

	updateErr := errors.New("update failed")

	cart := &models.CartItem{
		ID:       1,
		UserID:   1,
		Quantity: 2,
		Product: models.Product{
			ID:   10,
			Name: "Mechanical Keyboard",
			Category: models.Category{
				ID:   2,
				Name: "Gaming",
			},
		},
	}

	cartRepo.
		On("GetByID", uint(1)).
		Return(cart, nil)

	cartRepo.
		On("Update", mock.AnythingOfType("*models.CartItem")).
		Return(updateErr)

	resp, err := service.Update(1, 1, req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, updateErr)

	cartRepo.AssertExpectations(t)
}

func TestDeleteCartSuccess(t *testing.T) {
	cartRepo := new(serviceMocks.CartRepository)
	productRepo := new(serviceMocks.ProductRepository)

	service := services.NewCartService(cartRepo, productRepo)

	cartRepo.
		On("GetByID", uint(1)).
		Return(&models.CartItem{
			ID:     1,
			UserID: 1,
		}, nil)

	cartRepo.
		On("Delete", uint(1)).
		Return(nil)

	err := service.Delete(1, 1)

	assert.NoError(t, err)

	cartRepo.AssertExpectations(t)
}

func TestDeleteCartNotFound(t *testing.T) {
	cartRepo := new(serviceMocks.CartRepository)
	productRepo := new(serviceMocks.ProductRepository)

	service := services.NewCartService(cartRepo, productRepo)

	cartRepo.
		On("GetByID", uint(1)).
		Return(nil, nil)

	err := service.Delete(1, 1)

	assert.ErrorIs(t, err, errs.ErrCartItemNotFound)

	cartRepo.AssertNotCalled(t, "Delete", uint(1))

	cartRepo.AssertExpectations(t)
}

func TestDeleteCartForbidden(t *testing.T) {
	cartRepo := new(serviceMocks.CartRepository)
	productRepo := new(serviceMocks.ProductRepository)

	service := services.NewCartService(cartRepo, productRepo)

	cartRepo.
		On("GetByID", uint(1)).
		Return(&models.CartItem{
			ID:     1,
			UserID: 99,
		}, nil)

	err := service.Delete(1, 1)

	assert.ErrorIs(t, err, errs.ErrForbidden)

	cartRepo.AssertNotCalled(t, "Delete", uint(1))

	cartRepo.AssertExpectations(t)
}

func TestDeleteCartRepositoryError(t *testing.T) {
	cartRepo := new(serviceMocks.CartRepository)
	productRepo := new(serviceMocks.ProductRepository)

	service := services.NewCartService(cartRepo, productRepo)

	dbErr := errors.New("database error")

	cartRepo.
		On("GetByID", uint(1)).
		Return(nil, dbErr)

	err := service.Delete(1, 1)

	assert.ErrorIs(t, err, dbErr)

	cartRepo.AssertNotCalled(t, "Delete", uint(1))

	cartRepo.AssertExpectations(t)
}

func TestDeleteCartDeleteError(t *testing.T) {
	cartRepo := new(serviceMocks.CartRepository)
	productRepo := new(serviceMocks.ProductRepository)

	service := services.NewCartService(cartRepo, productRepo)

	deleteErr := errors.New("delete failed")

	cartRepo.
		On("GetByID", uint(1)).
		Return(&models.CartItem{
			ID:     1,
			UserID: 1,
		}, nil)

	cartRepo.
		On("Delete", uint(1)).
		Return(deleteErr)

	err := service.Delete(1, 1)

	assert.ErrorIs(t, err, deleteErr)

	cartRepo.AssertExpectations(t)
}
