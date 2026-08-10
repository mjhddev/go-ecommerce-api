package services

import (
	"github.com/mjhddev/go-ecommerce-api/internal/dto"
	"github.com/mjhddev/go-ecommerce-api/internal/errs"
	"github.com/mjhddev/go-ecommerce-api/internal/models"
	"github.com/mjhddev/go-ecommerce-api/internal/repositories"
)

type CartService interface {
	AddToCart(userID uint, request dto.AddToCartRequest) (*dto.CartItemResponse, error)
	GetCart(userID uint) ([]dto.CartItemResponse, error)
	Update(userID, cartID uint, request dto.UpdateCartRequest) (*dto.CartItemResponse, error)
	Delete(userID, cartID uint) error
}

type cartService struct {
	cartRepo    repositories.CartRepository
	productRepo repositories.ProductRepository
}

func NewCartService(
	cartRepo repositories.CartRepository,
	productRepo repositories.ProductRepository,
) CartService {
	return &cartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (s *cartService) AddToCart(userID uint, request dto.AddToCartRequest) (*dto.CartItemResponse, error) {
	// Check if the product exists
	product, err := s.productRepo.GetByID(request.ProductID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errs.ErrProductNotFound
	}

	// Check if the product is already in the cart
	cartItem, err := s.cartRepo.GetByUserAndProduct(userID, request.ProductID)
	if err != nil {
		return nil, err
	}

	if cartItem != nil {
		// If the product is already in the cart, update the quantity
		cartItem.Quantity += request.Quantity
		err = s.cartRepo.Update(cartItem)
		if err != nil {
			return nil, err
		}
		return &dto.CartItemResponse{
			ID:       cartItem.ID,
			Quantity: cartItem.Quantity,
			Product: dto.ProductResponse{
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
			},
		}, nil
	}

	cartItem = &models.CartItem{
		UserID:    userID,
		ProductID: request.ProductID,
		Quantity:  request.Quantity,
	}

	err = s.cartRepo.Create(cartItem)
	if err != nil {
		return nil, err
	}

	return &dto.CartItemResponse{
		ID:       cartItem.ID,
		Quantity: cartItem.Quantity,
		Product: dto.ProductResponse{
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
		},
	}, nil
}

func (s *cartService) GetCart(userID uint) ([]dto.CartItemResponse, error) {
	cartItems, err := s.cartRepo.GetByUser(userID)
	if err != nil {
		return nil, err
	}

	response := make([]dto.CartItemResponse, 0, len(cartItems))
	for _, item := range cartItems {
		response = append(response, dto.CartItemResponse{
			ID:       item.ID,
			Quantity: item.Quantity,
			Product: dto.ProductResponse{
				ID:          item.Product.ID,
				Name:        item.Product.Name,
				Description: item.Product.Description,
				Price:       item.Product.Price,
				Stock:       item.Product.Stock,
				ImageURL:    item.Product.ImageURL,

				Category: dto.CategoryResponse{
					ID:   item.Product.Category.ID,
					Name: item.Product.Category.Name,
				},
			},
		})
	}

	return response, nil
}

func (s *cartService) Update(userID, cartID uint, request dto.UpdateCartRequest) (*dto.CartItemResponse, error) {
	cartItem, err := s.cartRepo.GetByID(cartID)
	if err != nil {
		return nil, err
	}
	if cartItem == nil {
		return nil, errs.ErrCartItemNotFound
	}
	if cartItem.UserID != userID {
		return nil, errs.ErrForbidden
	}

	cartItem.Quantity = request.Quantity
	err = s.cartRepo.Update(cartItem)
	if err != nil {
		return nil, err
	}

	return &dto.CartItemResponse{
		ID:       cartItem.ID,
		Quantity: cartItem.Quantity,
		Product: dto.ProductResponse{
			ID:          cartItem.Product.ID,
			Name:        cartItem.Product.Name,
			Description: cartItem.Product.Description,
			Price:       cartItem.Product.Price,
			Stock:       cartItem.Product.Stock,
			ImageURL:    cartItem.Product.ImageURL,

			Category: dto.CategoryResponse{
				ID:   cartItem.Product.Category.ID,
				Name: cartItem.Product.Category.Name,
			},
		},
	}, nil
}

func (s *cartService) Delete(userID, cartID uint) error {
	cartItem, err := s.cartRepo.GetByID(cartID)
	if err != nil {
		return err
	}
	if cartItem == nil {
		return errs.ErrCartItemNotFound
	}
	if cartItem.UserID != userID {
		return errs.ErrForbidden
	}

	return s.cartRepo.Delete(cartID)
}
