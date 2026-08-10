package services

import (
	"math"
	"mime/multipart"

	"github.com/mjhddev/go-ecommerce-api/internal/dto"
	"github.com/mjhddev/go-ecommerce-api/internal/errs"
	"github.com/mjhddev/go-ecommerce-api/internal/models"
	"github.com/mjhddev/go-ecommerce-api/internal/repositories"
	"github.com/mjhddev/go-ecommerce-api/internal/storage"
)

type ProductService interface {
	Create(request dto.CreateProductRequest) (*dto.ProductResponse, error)
	GetAll(query dto.ProductQuery) (*dto.ProductListResponse, error)
	GetByID(id uint) (*dto.ProductResponse, error)
	Update(id uint, request dto.UpdateProductRequest) (*dto.ProductResponse, error)
	Delete(id uint) error

	UploadImage(id uint, file *multipart.FileHeader) (*dto.ProductResponse, error)
}

type productService struct {
	productRepo  repositories.ProductRepository
	categoryRepo repositories.CategoryRepository
}

func NewProductService(
	productRepo repositories.ProductRepository,
	categoryRepo repositories.CategoryRepository,
) ProductService {
	return &productService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *productService) Create(request dto.CreateProductRequest) (*dto.ProductResponse, error) {
	category, err := s.categoryRepo.GetByID(request.CategoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errs.ErrCategoryNotFound
	}

	product := &models.Product{
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Stock:       request.Stock,
		CategoryID:  request.CategoryID,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, err
	}

	response := &dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		ImageURL:    product.ImageURL,
		Category: dto.CategoryResponse{
			ID:   category.ID,
			Name: category.Name,
		},
	}

	return response, nil
}

func (s *productService) GetAll(query dto.ProductQuery) (*dto.ProductListResponse, error) {
	products, total, err := s.productRepo.GetAll(query)
	if err != nil {
		return nil, err
	}

	items := make([]dto.ProductResponse, 0, len(products))

	for _, product := range products {
		items = append(items, dto.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			ImageURL:    product.ImageURL,

			Category: dto.CategoryResponse{
				ID:   product.Category.ID,
				Name: product.Category.Name,
			},
		})
	}

	totalPage := int(math.Ceil(float64(total) / float64(query.Limit)))

	return &dto.ProductListResponse{
		Items: items,
		Pagination: dto.PaginationResponse{
			Page:      query.Page,
			Limit:     query.Limit,
			TotalData: total,
			TotalPage: totalPage,
		},
	}, nil
}

func (s *productService) GetByID(id uint) (*dto.ProductResponse, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errs.ErrProductNotFound
	}

	response := &dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		ImageURL:    product.ImageURL,

		Category: dto.CategoryResponse{
			ID:   product.Category.ID,
			Name: product.Category.Name,
		},
	}

	return response, nil
}

func (s *productService) Update(id uint, request dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errs.ErrProductNotFound
	}

	category, err := s.categoryRepo.GetByID(request.CategoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errs.ErrCategoryNotFound
	}

	product.Name = request.Name
	product.Description = request.Description
	product.Price = request.Price
	product.Stock = request.Stock
	product.CategoryID = request.CategoryID

	if err := s.productRepo.Update(product); err != nil {
		return nil, err
	}

	response := &dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		ImageURL:    product.ImageURL,

		Category: dto.CategoryResponse{
			ID:   category.ID,
			Name: category.Name,
		},
	}

	return response, nil
}

func (s *productService) Delete(id uint) error {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return err
	}
	if product == nil {
		return errs.ErrProductNotFound
	}

	return s.productRepo.Delete(id)
}

func (s *productService) UploadImage(id uint, file *multipart.FileHeader) (*dto.ProductResponse, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, errs.ErrProductNotFound
	}
	oldImage := product.ImageURL

	imageURL, err := storage.SaveProductImage(file)
	if err != nil {
		return nil, err
	}

	product.ImageURL = imageURL

	if err := s.productRepo.Update(product); err != nil {
		// Database gagal, hapus file yang baru diupload agar tidak menjadi orphan
		_ = storage.DeleteFile(imageURL)
		return nil, err
	}
	// Database berhasil, sekarang hapus gambar lama
	if oldImage != "" {
		_ = storage.DeleteFile(oldImage)
	}

	return &dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		ImageURL:    product.ImageURL,
		Category: dto.CategoryResponse{
			ID:   product.Category.ID,
			Name: product.Category.Name,
		},
	}, nil
}
