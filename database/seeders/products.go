package seeders

import (
	"github.com/mjhddev/go-ecommerce-api/internal/models"
	"gorm.io/gorm"
)

func seedProducts(db *gorm.DB) error {
	products := []models.Product{
		{
			Name:        "Gaming Laptop",
			Description: "15.6 inch Gaming Laptop",
			Price:       12500000,
			Stock:       7,
			CategoryID:  1,
		},
		{
			Name:        "Bluetooth Speaker",
			Description: "Portable Bluetooth Speaker",
			Price:       550000,
			Stock:       20,
			CategoryID:  1,
		},
		{
			Name:        "Smart Watch",
			Description: "Fitness Smart Watch",
			Price:       950000,
			Stock:       14,
			CategoryID:  1,
		},
		{
			Name:        "Power Bank 20000mAh",
			Description: "Fast Charging Power Bank",
			Price:       420000,
			Stock:       35,
			CategoryID:  1,
		},
		{
			Name:        "Wireless Charger",
			Description: "15W Wireless Charger",
			Price:       280000,
			Stock:       18,
			CategoryID:  1,
		},
		{
			Name:        "Gaming Controller",
			Description: "Wireless Game Controller",
			Price:       620000,
			Stock:       16,
			CategoryID:  2,
		},
		{
			Name:        "Mouse Pad XL",
			Description: "Large RGB Mouse Pad",
			Price:       175000,
			Stock:       50,
			CategoryID:  2,
		},
		{
			Name:        "Gaming Desk",
			Description: "Carbon Fiber Gaming Desk",
			Price:       2100000,
			Stock:       5,
			CategoryID:  2,
		},
		{
			Name:        "VR Headset",
			Description: "Virtual Reality Headset",
			Price:       4200000,
			Stock:       4,
			CategoryID:  2,
		},
		{
			Name:        "Gaming Monitor 32 Inch",
			Description: "165Hz Gaming Monitor",
			Price:       3500000,
			Stock:       9,
			CategoryID:  2,
		},
		{
			Name:        "Men T-Shirt",
			Description: "Cotton T-Shirt",
			Price:       120000,
			Stock:       60,
			CategoryID:  3,
		},
		{
			Name:        "Hoodie",
			Description: "Oversized Hoodie",
			Price:       280000,
			Stock:       25,
			CategoryID:  3,
		},
		{
			Name:        "Running Shoes",
			Description: "Lightweight Running Shoes",
			Price:       850000,
			Stock:       18,
			CategoryID:  3,
		},
		{
			Name:        "Baseball Cap",
			Description: "Adjustable Cap",
			Price:       95000,
			Stock:       45,
			CategoryID:  3,
		},
		{
			Name:        "Backpack",
			Description: "Waterproof Backpack",
			Price:       390000,
			Stock:       21,
			CategoryID:  3,
		},
		{
			Name:        "Clean Code",
			Description: "Book by Robert C. Martin",
			Price:       320000,
			Stock:       12,
			CategoryID:  4,
		},
		{
			Name:        "The Pragmatic Programmer",
			Description: "Programming Best Practices",
			Price:       410000,
			Stock:       10,
			CategoryID:  4,
		},
		{
			Name:        "Design Patterns",
			Description: "Elements of Reusable Object-Oriented Software",
			Price:       550000,
			Stock:       8,
			CategoryID:  4,
		},
		{
			Name:        "Go Programming",
			Description: "Learning Go Language",
			Price:       295000,
			Stock:       15,
			CategoryID:  4,
		},
		{
			Name:        "Algorithms Book",
			Description: "Data Structures and Algorithms",
			Price:       470000,
			Stock:       9,
			CategoryID:  4,
		},
		{
			Name:        "Arabica Coffee Beans",
			Description: "Premium Coffee Beans 1kg",
			Price:       185000,
			Stock:       32,
			CategoryID:  5,
		},
		{
			Name:        "Green Tea",
			Description: "Organic Green Tea",
			Price:       75000,
			Stock:       28,
			CategoryID:  5,
		},
		{
			Name:        "Dark Chocolate",
			Description: "70% Cocoa Chocolate",
			Price:       45000,
			Stock:       40,
			CategoryID:  5,
		},
		{
			Name:        "Instant Noodles Box",
			Description: "Box of 40 Packs",
			Price:       135000,
			Stock:       17,
			CategoryID:  5,
		},
		{
			Name:        "Potato Chips",
			Description: "Original Flavor",
			Price:       25000,
			Stock:       80,
			CategoryID:  5,
		},
	}

	for _, product := range products {
		var existing models.Product

		if err := db.Where("name = ?", product.Name).First(&existing).Error; err == nil {
			continue
		}

		if err := db.Create(&product).Error; err != nil {
			return err
		}
	}

	return nil
}
