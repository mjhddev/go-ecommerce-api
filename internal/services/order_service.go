package services

import (
	"github.com/mjhddev/go-ecommerce-api/internal/dto"
	"github.com/mjhddev/go-ecommerce-api/internal/errs"
	"github.com/mjhddev/go-ecommerce-api/internal/models"
	"github.com/mjhddev/go-ecommerce-api/internal/repositories"
	"gorm.io/gorm"
)

type OrderService interface {
	Checkout(userID uint) (*dto.CheckoutResponse, error)
}

type orderService struct {
	db          *gorm.DB
	cartRepo    repositories.CartRepository
	productRepo repositories.ProductRepository
	orderRepo   repositories.OrderRepository
}

func NewOrderService(
	db *gorm.DB,
	orderRepo repositories.OrderRepository,
	cartRepo repositories.CartRepository,
	productRepo repositories.ProductRepository,
) OrderService {
	return &orderService{
		db:          db,
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (s *orderService) Checkout(userID uint) (*dto.CheckoutResponse, error) {
	cartItems, err := s.cartRepo.GetByUser(userID)
	if err != nil {
		return nil, err
	}

	if len(cartItems) == 0 {
		return nil, errs.ErrCartEmpty
	}

	var total float64
	for _, item := range cartItems {
		total += item.Product.Price * float64(item.Quantity)
	}

	order := &models.Order{
		UserID:      userID,
		TotalAmount: total,
		Status:      "pending",
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.orderRepo.Checkout(tx, order); err != nil {
			return err
		}

		for _, item := range cartItems {

			if item.Product.Stock < item.Quantity {
				return errs.ErrInsufficientStock
			}

			orderItem := &models.OrderItem{
				OrderID:   order.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Product.Price,
			}
			if err := s.orderRepo.CreateOrderItem(tx, orderItem); err != nil {
				return err
			}

			item.Product.Stock -= item.Quantity
			if err := s.productRepo.UpdateTx(tx, &item.Product); err != nil {
				return err
			}
		}

		if err := s.cartRepo.DeleteByUser(tx, userID); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &dto.CheckoutResponse{
		OrderID:     order.ID,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
	}, nil
}
