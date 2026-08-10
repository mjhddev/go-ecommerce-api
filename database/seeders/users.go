package seeders

import (
	"github.com/mjhddev/go-ecommerce-api/internal/models"
	"github.com/mjhddev/go-ecommerce-api/internal/utils"
	"gorm.io/gorm"
)

func seedUsers(db *gorm.DB) error {
	users := []models.User{
		{
			Name:  "Administrator",
			Email: "admin@example.com",
			Role:  "admin",
		},
		{
			Name:  "John Doe",
			Email: "john@example.com",
			Role:  "user",
		},
		{
			Name:  "Jane Doe",
			Email: "jane@example.com",
			Role:  "user",
		},
		{
			Name:  "Bob Smith",
			Email: "bob@example.com",
			Role:  "user",
		},
	}

	for _, user := range users {
		var existing models.User

		if err := db.Where("email = ?", user.Email).First(&existing).Error; err == nil {
			continue
		}

		password, err := utils.HashPassword("password123")
		if err != nil {
			return err
		}

		user.Password = password

		if err := db.Create(&user).Error; err != nil {
			return err
		}
	}

	return nil
}
