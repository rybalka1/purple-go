package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rybalka1/purple-go/4-order-api/internal/config"
	"github.com/rybalka1/purple-go/4-order-api/internal/database"
	"github.com/rybalka1/purple-go/4-order-api/internal/handlers"
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

	r := mux.NewRouter()

	r.HandleFunc("/product", handlers.CreateProduct).Methods("POST")
	r.HandleFunc("/product/{id}", handlers.GetProduct).Methods("GET")
	r.HandleFunc("/product/{id}", handlers.UpdateProduct).Methods("PUT")
	r.HandleFunc("/product/{id}", handlers.DeleteProduct).Methods("DELETE")

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
