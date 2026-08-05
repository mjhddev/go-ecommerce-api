package models

import "time"

type CartItem struct {
	ID uint `gorm:"primaryKey" json:"id"`

	UserID uint `json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	ProductID uint    `json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`

	Quantity int `gorm:"not null;default:1" json:"quantity"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
