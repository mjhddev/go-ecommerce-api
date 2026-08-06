package main

import (
	"log"
	"os"

	_ "github.com/mjhddev/go-ecommerce-api/docs"

	"github.com/mjhddev/go-ecommerce-api/configs"
	"github.com/mjhddev/go-ecommerce-api/internal/routes"
)

// @title Go E-Commerce API
// @version 1.0
// @description RESTful E-Commerce API built with Go, Gin, GORM and PostgreSQL.
// @BasePath /api/v1
// @host localhost:8080
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	log.Println("Loading environment...")
	configs.LoadEnv()

	log.Println("Connecting database...")
	configs.ConnectDatabase()

	log.Println("Setting up router...")
	router := routes.SetupRouter()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)

	router.Run(":" + port)
}
