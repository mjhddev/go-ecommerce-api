package models

import "time"

type Order struct {
	ID uint `gorm:"primaryKey" json:"id"`

	UserID uint `json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	TotalAmount float64 `gorm:"not null" json:"total_amount"`

	Status string `gorm:"default:'pending'" json:"status"`

	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"order_items,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
