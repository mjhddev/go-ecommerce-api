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

func TestCreateProductSuccess(t *testing.T) {

	productRepo := new(serviceMocks.ProductRepository)
	categoryRepo := new(serviceMocks.CategoryRepository)

	service := services.NewProductService(productRepo, categoryRepo)

	req := dto.CreateProductRequest{
		Name:        "Mechanical Keyboard",
		Description: "RGB Keyboard",
		Price:       750000,
		Stock:       20,
		CategoryID:  2,
	}

	categoryRepo.
		On("GetByID", req.CategoryID).
		Return(&models.Category{
			ID:   2,
			Name: "Gaming",
		}, nil)

	productRepo.
		On("Create", mock.AnythingOfType("*models.Product")).
		Run(func(args mock.Arguments) {

			product := args.Get(0).(*models.Product)

			product.ID = 10

		}).
		Return(nil)

	resp, err := service.Create(req)

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, uint(10), resp.ID)
	assert.Equal(t, req.Name, resp.Name)
	assert.Equal(t, req.Description, resp.Description)
	assert.Equal(t, req.Price, resp.Price)
	assert.Equal(t, req.Stock, resp.Stock)

	assert.Equal(t, uint(2), resp.Category.ID)
	assert.Equal(t, "Gaming", resp.Category.Name)

	productRepo.AssertExpectations(t)
	categoryRepo.AssertExpectations(t)
}

func TestCreateProductCategoryNotFound(t *testing.T) {
	productRepo := new(serviceMocks.ProductRepository)
	categoryRepo := new(serviceMocks.CategoryRepository)

	service := services.NewProductService(productRepo, categoryRepo)

	req := dto.CreateProductRequest{
		Name:        "Keyboard",
		Description: "RGB",
		Price:       750000,
		Stock:       10,
		CategoryID:  100,
	}

	categoryRepo.
		On("GetByID", req.CategoryID).
		Return(nil, nil)

	resp, err := service.Create(req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, errs.ErrCategoryNotFound)

	productRepo.AssertNotCalled(t, "Create", mock.Anything)

	categoryRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}

func TestCreateProductCategoryRepositoryError(t *testing.T) {
	productRepo := new(serviceMocks.ProductRepository)
	categoryRepo := new(serviceMocks.CategoryRepository)

	service := services.NewProductService(productRepo, categoryRepo)

	req := dto.CreateProductRequest{
		CategoryID: 2,
	}

	dbErr := errors.New("database error")

	categoryRepo.
		On("GetByID", req.CategoryID).
		Return(nil, dbErr)

	resp, err := service.Create(req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, dbErr)

	productRepo.AssertNotCalled(t, "Create", mock.Anything)

	categoryRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}

func TestCreateProductCreateError(t *testing.T) {
	productRepo := new(serviceMocks.ProductRepository)
	categoryRepo := new(serviceMocks.CategoryRepository)

	service := services.NewProductService(productRepo, categoryRepo)

	req := dto.CreateProductRequest{
		Name:        "Keyboard",
		Description: "RGB",
		Price:       750000,
		Stock:       10,
		CategoryID:  2,
	}

	createErr := errors.New("create failed")

	categoryRepo.
		On("GetByID", req.CategoryID).
		Return(&models.Category{
			ID:   2,
			Name: "Gaming",
		}, nil)

	productRepo.
		On("Create", mock.AnythingOfType("*models.Product")).
		Return(createErr)

	resp, err := service.Create(req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, createErr)

	productRepo.AssertExpectations(t)
	categoryRepo.AssertExpectations(t)
}

func TestGetAllProductsSuccess(t *testing.T) {
	productRepo := new(serviceMocks.ProductRepository)
	categoryRepo := new(serviceMocks.CategoryRepository)

	service := services.NewProductService(productRepo, categoryRepo)

	products := []models.Product{
		{
			ID:          1,
			Name:        "Keyboard",
			Description: "RGB",
			Price:       750000,
			Stock:       20,
			Category: models.Category{
				ID:   2,
				Name: "Gaming",
			},
		},
		{
			ID:          2,
			Name:        "Mouse",
			Description: "Wireless",
			Price:       250000,
			Stock:       15,
			Category: models.Category{
				ID:   2,
				Name: "Gaming",
			},
		},
	}

	productRepo.
		On("GetAll").
		Return(products, nil)

	resp, err := service.GetAll()

	require.NoError(t, err)
	require.Len(t, resp, 2)

	assert.Equal(t, "Keyboard", resp[0].Name)
	assert.Equal(t, "Gaming", resp[0].Category.Name)

	productRepo.AssertExpectations(t)
}

func TestGetAllProductsError(t *testing.T) {
	productRepo := new(serviceMocks.ProductRepository)
	categoryRepo := new(serviceMocks.CategoryRepository)

	service := services.NewProductService(productRepo, categoryRepo)

	dbErr := errors.New("database error")

	productRepo.
		On("GetAll").
		Return(nil, dbErr)

	resp, err := service.GetAll()

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, dbErr)

	productRepo.AssertExpectations(t)
}

func TestGetProductByIDSuccess(t *testing.T) {
	productRepo := new(serviceMocks.ProductRepository)
	categoryRepo := new(serviceMocks.CategoryRepository)

	service := services.NewProductService(productRepo, categoryRepo)

	product := &models.Product{
		ID:          1,
		Name:        "Keyboard",
		Description: "RGB",
		Price:       750000,
		Stock:       20,
		Category: models.Category{
			ID:   2,
			Name: "Gaming",
		},
	}

	productRepo.
		On("GetByID", uint(1)).
		Return(product, nil)

	resp, err := service.GetByID(1)

	require.NoError(t, err)

	assert.Equal(t, product.Name, resp.Name)
	assert.Equal(t, product.Category.Name, resp.Category.Name)

	productRepo.AssertExpectations(t)
}

func TestGetProductByIDNotFound(t *testing.T) {
	productRepo := new(serviceMocks.ProductRepository)
	categoryRepo := new(serviceMocks.CategoryRepository)

	service := services.NewProductService(productRepo, categoryRepo)

	productRepo.
		On("GetByID", uint(1)).
		Return(nil, nil)

	resp, err := service.GetByID(1)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, errs.ErrProductNotFound)

	productRepo.AssertExpectations(t)
}

func TestGetProductByIDRepositoryError(t *testing.T) {
	productRepo := new(serviceMocks.ProductRepository)
	categoryRepo := new(serviceMocks.CategoryRepository)

	service := services.NewProductService(productRepo, categoryRepo)

	dbErr := errors.New("database error")

	productRepo.
		On("GetByID", uint(1)).
		Return(nil, dbErr)

	resp, err := service.GetByID(1)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, dbErr)

	productRepo.AssertExpectations(t)
}

func TestUpdateProductSuccess(t *testing.T) {
	productRepo := new(serviceMocks.ProductRepository)
	categoryRepo := new(serviceMocks.CategoryRepository)

	service := services.NewProductService(productRepo, categoryRepo)

	req := dto.UpdateProductRequest{
		Name:        "Mechanical Keyboard",
		Description: "RGB Keyboard",
		Price:       800000,
		Stock:       25,
		CategoryID:  2,
	}

	product := &models.Product{
		ID:          1,
		Name:        "Old Keyboard",
		Description: "Old",
		Price:       700000,
		Stock:       10,
		CategoryID:  1,
	}

	category := &models.Category{
		ID:   2,
		Name: "Gaming",
	}

	productRepo.
		On("GetByID", uint(1)).
		Return(product, nil)

	categoryRepo.
		On("GetByID", uint(2)).
		Return(category, nil)

	productRepo.
		On("Update", mock.AnythingOfType("*models.Product")).
		Run(func(args mock.Arguments) {
			p := args.Get(0).(*models.Product)

			assert.Equal(t, req.Name, p.Name)
			assert.Equal(t, req.Description, p.Description)
			assert.Equal(t, req.Price, p.Price)
			assert.Equal(t, req.Stock, p.Stock)
			assert.Equal(t, req.CategoryID, p.CategoryID)
		}).
		Return(nil)

	resp, err := service.Update(1, req)

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, req.Name, resp.Name)
	assert.Equal(t, req.Description, resp.Description)
	assert.Equal(t, req.Price, resp.Price)
	assert.Equal(t, req.Stock, resp.Stock)
	assert.Equal(t, category.Name, resp.Category.Name)

	productRepo.AssertExpectations(t)
	categoryRepo.AssertExpectations(t)
}

func TestUpdateProductNotFound(t *testing.T) {
	productRepo := new(serviceMocks.ProductRepository)
	categoryRepo := new(serviceMocks.CategoryRepository)

	service := services.NewProductService(productRepo, categoryRepo)

	req := dto.UpdateProductRequest{}

	productRepo.
		On("GetByID", uint(1)).
		Return(nil, nil)

	resp, err := service.Update(1, req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, errs.ErrProductNotFound)

	categoryRepo.AssertNotCalled(t, "GetByID", mock.Anything)
	productRepo.AssertNotCalled(t, "Update", mock.Anything)

	productRepo.AssertExpectations(t)
}

func TestUpdateProductCategoryNotFound(t *testing.T) {
	productRepo := new(serviceMocks.ProductRepository)
	categoryRepo := new(serviceMocks.CategoryRepository)

	service := services.NewProductService(productRepo, categoryRepo)

	req := dto.UpdateProductRequest{
		CategoryID: 99,
	}

	productRepo.
		On("GetByID", uint(1)).
		Return(&models.Product{ID: 1}, nil)

	categoryRepo.
		On("GetByID", uint(99)).
		Return(nil, nil)

	resp, err := service.Update(1, req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, errs.ErrCategoryNotFound)

	productRepo.AssertNotCalled(t, "Update", mock.Anything)

	productRepo.AssertExpectations(t)
	categoryRepo.AssertExpectations(t)
}

func TestUpdateProductUpdateError(t *testing.T) {
	productRepo := new(serviceMocks.ProductRepository)
	categoryRepo := new(serviceMocks.CategoryRepository)

	service := services.NewProductService(productRepo, categoryRepo)

	req := dto.UpdateProductRequest{
		Name:        "Keyboard",
		Description: "RGB",
		Price:       750000,
		Stock:       20,
		CategoryID:  2,
	}

	updateErr := errors.New("update failed")

	productRepo.
		On("GetByID", uint(1)).
		Return(&models.Product{ID: 1}, nil)

	categoryRepo.
		On("GetByID", uint(2)).
		Return(&models.Category{
			ID:   2,
			Name: "Gaming",
		}, nil)

	productRepo.
		On("Update", mock.AnythingOfType("*models.Product")).
		Return(updateErr)

	resp, err := service.Update(1, req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, updateErr)

	productRepo.AssertExpectations(t)
	categoryRepo.AssertExpectations(t)
}

func TestDeleteProductSuccess(t *testing.T) {
	productRepo := new(serviceMocks.ProductRepository)
	categoryRepo := new(serviceMocks.CategoryRepository)

	service := services.NewProductService(productRepo, categoryRepo)

	productRepo.
		On("GetByID", uint(1)).
		Return(&models.Product{ID: 1}, nil)

	productRepo.
		On("Delete", uint(1)).
		Return(nil)

	err := service.Delete(1)

	assert.NoError(t, err)

	productRepo.AssertExpectations(t)
}

func TestDeleteProductNotFound(t *testing.T) {
	productRepo := new(serviceMocks.ProductRepository)
	categoryRepo := new(serviceMocks.CategoryRepository)

	service := services.NewProductService(productRepo, categoryRepo)

	productRepo.
		On("GetByID", uint(1)).
		Return(nil, nil)

	err := service.Delete(1)

	assert.ErrorIs(t, err, errs.ErrProductNotFound)

	productRepo.AssertNotCalled(t, "Delete", uint(1))

	productRepo.AssertExpectations(t)
}

func TestDeleteProductRepositoryError(t *testing.T) {
	productRepo := new(serviceMocks.ProductRepository)
	categoryRepo := new(serviceMocks.CategoryRepository)

	service := services.NewProductService(productRepo, categoryRepo)

	dbErr := errors.New("database error")

	productRepo.
		On("GetByID", uint(1)).
		Return(nil, dbErr)

	err := service.Delete(1)

	assert.ErrorIs(t, err, dbErr)

	productRepo.AssertNotCalled(t, "Delete", uint(1))

	productRepo.AssertExpectations(t)
}

func TestDeleteProductDeleteError(t *testing.T) {
	productRepo := new(serviceMocks.ProductRepository)
	categoryRepo := new(serviceMocks.CategoryRepository)

	service := services.NewProductService(productRepo, categoryRepo)

	deleteErr := errors.New("delete failed")

	productRepo.
		On("GetByID", uint(1)).
		Return(&models.Product{ID: 1}, nil)

	productRepo.
		On("Delete", uint(1)).
		Return(deleteErr)

	err := service.Delete(1)

	assert.ErrorIs(t, err, deleteErr)

	productRepo.AssertExpectations(t)
}
