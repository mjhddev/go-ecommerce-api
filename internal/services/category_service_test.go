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

func TestCreateCategorySuccess(t *testing.T) {
	mockRepo := new(serviceMocks.CategoryRepository)

	service := services.NewCategoryService(mockRepo)

	req := dto.CreateCategoryRequest{
		Name: "Gaming",
	}

	mockRepo.
		On("Create", mock.AnythingOfType("*models.Category")).
		Run(func(args mock.Arguments) {
			category := args.Get(0).(*models.Category)

			category.ID = 1
		}).
		Return(nil)

	resp, err := service.Create(req)

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, uint(1), resp.ID)
	assert.Equal(t, req.Name, resp.Name)

	mockRepo.AssertExpectations(t)
}

func TestCreateCategoryError(t *testing.T) {
	mockRepo := new(serviceMocks.CategoryRepository)

	service := services.NewCategoryService(mockRepo)

	req := dto.CreateCategoryRequest{
		Name: "Gaming",
	}

	createErr := errors.New("create failed")

	mockRepo.
		On("Create", mock.AnythingOfType("*models.Category")).
		Return(createErr)

	resp, err := service.Create(req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.ErrorIs(t, err, createErr)

	mockRepo.AssertExpectations(t)
}

func TestGetAllCategories(t *testing.T) {
	tests := []struct {
		name      string
		mockData  []models.Category
		mockErr   error
		wantCount int
		wantErr   bool
	}{
		{
			name: "Success",
			mockData: []models.Category{
				{
					ID:   1,
					Name: "Gaming",
				},
				{
					ID:   2,
					Name: "Office",
				},
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "Repository Error",
			mockErr:   errors.New("database error"),
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := new(serviceMocks.CategoryRepository)

			service := services.NewCategoryService(mockRepo)

			mockRepo.
				On("GetAll").
				Return(tt.mockData, tt.mockErr)

			resp, err := service.GetAll()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Len(t, resp, tt.wantCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetCategoryByID(t *testing.T) {
	tests := []struct {
		name     string
		id       uint
		mockData *models.Category
		mockErr  error
		wantErr  error
		wantName string
	}{
		{
			name: "Success",
			id:   1,
			mockData: &models.Category{
				ID:   1,
				Name: "Gaming",
			},
			wantName: "Gaming",
		},
		{
			name:     "Category Not Found",
			id:       99,
			mockData: nil,
			wantErr:  errs.ErrCategoryNotFound,
		},
		{
			name:    "Repository Error",
			id:      1,
			mockErr: errors.New("database error"),
			wantErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := new(serviceMocks.CategoryRepository)

			service := services.NewCategoryService(mockRepo)

			mockRepo.
				On("GetByID", tt.id).
				Return(tt.mockData, tt.mockErr)

			resp, err := service.GetByID(tt.id)

			if tt.wantErr != nil {
				assert.Nil(t, resp)
				assert.Error(t, err)

				// Khusus business error
				if errors.Is(tt.wantErr, errs.ErrCategoryNotFound) {
					assert.ErrorIs(t, err, errs.ErrCategoryNotFound)
				} else {
					assert.EqualError(t, err, tt.wantErr.Error())
				}

			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)

				assert.Equal(t, tt.mockData.ID, resp.ID)
				assert.Equal(t, tt.wantName, resp.Name)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateCategorySuccess(t *testing.T) {
	mockRepo := new(serviceMocks.CategoryRepository)
	service := services.NewCategoryService(mockRepo)

	req := dto.UpdateCategoryRequest{
		Name: "Updated Gaming",
	}

	category := &models.Category{
		ID:   1,
		Name: "Gaming",
	}

	mockRepo.
		On("GetByID", uint(1)).
		Return(category, nil)

	mockRepo.
		On("Update", mock.AnythingOfType("*models.Category")).
		Run(func(args mock.Arguments) {
			c := args.Get(0).(*models.Category)
			assert.Equal(t, "Updated Gaming", c.Name)
		}).
		Return(nil)

	resp, err := service.Update(1, req)

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, uint(1), resp.ID)
	assert.Equal(t, "Updated Gaming", resp.Name)

	mockRepo.AssertExpectations(t)
}

func TestUpdateCategoryNotFound(t *testing.T) {
	mockRepo := new(serviceMocks.CategoryRepository)
	service := services.NewCategoryService(mockRepo)

	req := dto.UpdateCategoryRequest{
		Name: "Updated Gaming",
	}

	mockRepo.
		On("GetByID", uint(1)).
		Return(nil, nil)

	resp, err := service.Update(1, req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, errs.ErrCategoryNotFound)

	mockRepo.AssertNotCalled(t, "Update", mock.Anything)

	mockRepo.AssertExpectations(t)
}

func TestUpdateCategoryRepositoryError(t *testing.T) {
	mockRepo := new(serviceMocks.CategoryRepository)
	service := services.NewCategoryService(mockRepo)

	req := dto.UpdateCategoryRequest{
		Name: "Updated Gaming",
	}

	dbErr := errors.New("database error")

	mockRepo.
		On("GetByID", uint(1)).
		Return(nil, dbErr)

	resp, err := service.Update(1, req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, dbErr)

	mockRepo.AssertNotCalled(t, "Update", mock.Anything)

	mockRepo.AssertExpectations(t)
}

func TestUpdateCategoryUpdateError(t *testing.T) {
	mockRepo := new(serviceMocks.CategoryRepository)
	service := services.NewCategoryService(mockRepo)

	req := dto.UpdateCategoryRequest{
		Name: "Updated Gaming",
	}

	category := &models.Category{
		ID:   1,
		Name: "Gaming",
	}

	updateErr := errors.New("update failed")

	mockRepo.
		On("GetByID", uint(1)).
		Return(category, nil)

	mockRepo.
		On("Update", mock.AnythingOfType("*models.Category")).
		Return(updateErr)

	resp, err := service.Update(1, req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, updateErr)

	mockRepo.AssertExpectations(t)
}

func TestDeleteCategorySuccess(t *testing.T) {
	mockRepo := new(serviceMocks.CategoryRepository)
	service := services.NewCategoryService(mockRepo)

	mockRepo.
		On("GetByID", uint(1)).
		Return(&models.Category{
			ID: 1,
		}, nil)

	mockRepo.
		On("Delete", uint(1)).
		Return(nil)

	err := service.Delete(1)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestDeleteCategoryNotFound(t *testing.T) {
	mockRepo := new(serviceMocks.CategoryRepository)
	service := services.NewCategoryService(mockRepo)

	mockRepo.
		On("GetByID", uint(1)).
		Return(nil, nil)

	err := service.Delete(1)

	assert.ErrorIs(t, err, errs.ErrCategoryNotFound)

	mockRepo.AssertNotCalled(t, "Delete", uint(1))

	mockRepo.AssertExpectations(t)
}

func TestDeleteCategoryRepositoryError(t *testing.T) {
	mockRepo := new(serviceMocks.CategoryRepository)
	service := services.NewCategoryService(mockRepo)

	dbErr := errors.New("database error")

	mockRepo.
		On("GetByID", uint(1)).
		Return(nil, dbErr)

	err := service.Delete(1)

	assert.ErrorIs(t, err, dbErr)

	mockRepo.AssertNotCalled(t, "Delete", uint(1))

	mockRepo.AssertExpectations(t)
}

func TestDeleteCategoryDeleteError(t *testing.T) {
	mockRepo := new(serviceMocks.CategoryRepository)
	service := services.NewCategoryService(mockRepo)

	deleteErr := errors.New("delete failed")

	mockRepo.
		On("GetByID", uint(1)).
		Return(&models.Category{
			ID: 1,
		}, nil)

	mockRepo.
		On("Delete", uint(1)).
		Return(deleteErr)

	err := service.Delete(1)

	assert.ErrorIs(t, err, deleteErr)

	mockRepo.AssertExpectations(t)
}
