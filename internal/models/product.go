package models

import "time"

type Product struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string  `gorm:"size:255;not null" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	Price       float64 `gorm:"not null" json:"price"`
	Stock       int     `gorm:"default:0" json:"stock"`
	ImageURL    string  `gorm:"size:255" json:"image_url"`

	CategoryID uint     `json:"category_id"`
	Category   Category `gorm:"foreignKey:CategoryID" json:"category"`

	CartItems  []CartItem  `gorm:"foreignKey:ProductID" json:"cart_items,omitempty"`
	OrderItems []OrderItem `gorm:"foreignKey:ProductID" json:"order_items,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
