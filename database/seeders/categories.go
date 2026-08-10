package seeders

import (
	"github.com/mjhddev/go-ecommerce-api/internal/models"
	"gorm.io/gorm"
)

func seedCategories(db *gorm.DB) error {
	categories := []models.Category{
		{Name: "Electronics"},
		{Name: "Gaming"},
		{Name: "Fashion"},
		{Name: "Books"},
		{Name: "Food"},
	}

	for _, category := range categories {
		var existing models.Category

		if err := db.Where("name = ?", category.Name).First(&existing).Error; err == nil {
			continue
		}

		if err := db.Create(&category).Error; err != nil {
			return err
		}
	}

	return nil
}
