package main

import (
	"log"
	"os"

	"github.com/mjhddev/go-ecommerce-api/configs"
	"github.com/mjhddev/go-ecommerce-api/internal/routes"
)

func main() {
	configs.LoadEnv()
	configs.ConnectDatabase()

	router := routes.SetupRouter()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)

	router.Run(":" + port)
}
