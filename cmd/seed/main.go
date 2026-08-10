package main

import (
	"log"

	"github.com/mjhddev/go-ecommerce-api/configs"
	"github.com/mjhddev/go-ecommerce-api/database/seeders"
)

func main() {
	log.Println("Loading environment...")
	configs.LoadEnv()

	log.Println("Connecting database...")
	configs.ConnectDatabase()

	log.Println("Running database seeder...")

	if err := seeders.Seed(configs.DB); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}

	log.Println("Database seeded successfully")
}
