package models

import "time"

type Product struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string  `gorm:"size:255;not null" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	Price       float64 `gorm:"not null" json:"price"`
	Stock       int     `gorm:"default:0" json:"stock"`

	CategoryID uint     `json:"category_id"`
	Category   Category `gorm:"foreignKey:CategoryID" json:"category"`

	CartItems  []CartItem  `gorm:"foreignKey:ProductID" json:"cart_items,omitempty"`
	OrderItems []OrderItem `gorm:"foreignKey:ProductID" json:"order_items,omitempty"`

	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
}
