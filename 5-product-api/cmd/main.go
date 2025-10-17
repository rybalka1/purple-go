package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"purple-go/5-product-api/internal/config"
	"purple-go/5-product-api/internal/database"
	"purple-go/5-product-api/internal/handlers"
	"purple-go/5-product-api/internal/middleware"
	"purple-go/5-product-api/internal/models"
)

func main() {
	// Configure logrus to output in JSON format
	logrus.SetFormatter(&logrus.JSONFormatter{})

	cfg := config.Load()
	database.Connect(cfg)

	err := database.DB.AutoMigrate(&models.Product{})
	if err != nil {
		logrus.Fatalf("migration failed: %v", err)
	}

	logrus.Println("✅ Migrations complete. Product table created.")

	r := mux.NewRouter()

	// Apply the logging middleware
	r.Use(middleware.Logger)

	r.HandleFunc("/product", handlers.CreateProduct).Methods("POST")
	r.HandleFunc("/product/{id}", handlers.GetProduct).Methods("GET")
	r.HandleFunc("/product/{id}", handlers.UpdateProduct).Methods("PUT")
	r.HandleFunc("/product/{id}", handlers.DeleteProduct).Methods("DELETE")

	logrus.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		logrus.Fatal(err)
	}
}
