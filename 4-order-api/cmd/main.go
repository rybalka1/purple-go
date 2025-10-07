package main

import (
	"log"

	"github.com/rybalka1/purple-go/4-order-api/internal/config"
	"github.com/rybalka1/purple-go/4-order-api/internal/database"
	"github.com/rybalka1/purple-go/4-order-api/internal/models"
)

func main() {
	cfg := config.Load()
	database.Connect(cfg)

	err := database.DB.AutoMigrate(&models.Product{})
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("✅ Migrations complete. Product table created.")
}
