package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/mjhddev/go-ecommerce-api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB        *gorm.DB
	JWTSecret string
)

func ConnectDatabase() {
	JWTSecret = os.Getenv("JWT_SECRET")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to databse: ", err)
	}

	log.Println("Database connected successfully")

	err = DB.AutoMigrate(&models.User{})
}
