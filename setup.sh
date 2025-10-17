#!/bin/bash

PROJECT_DIR="/home/rybalka/Myprojects/purple-go/5-product-api"
SOURCE_DIR="/home/rybalka/Myprojects/purple-go/4-order-api"

# Create project directory
mkdir -p "$PROJECT_DIR"
cd "$PROJECT_DIR"

# Initialize Go module
go mod init purple-go/5-product-api

# Copy cmd and internal directories
cp -r "$SOURCE_DIR/cmd" ./
cp -r "$SOURCE_DIR/internal" ./

# Replace import paths in .go files
find . -type f -name "*.go" -print0 | xargs -0 sed -i 's|github.com/rybalka1/purple-go/4-order-api/internal|purple-go/5-product-api/internal|g'

# Add logrus dependency
go get github.com/sirupsen/logrus

# Create internal/middleware directory
mkdir -p internal/middleware

# Create internal/middleware/logger.go
cat << EOF > internal/middleware/logger.go
package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger returns a new http.Handler that logs requests to the given logger.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Our custom ResponseWriter to capture status code
		lw := &loggingResponseWriter{w, http.StatusOK}

		next.ServeHTTP(lw, r)

		duration := time.Since(start)

		logrus.WithFields(logrus.Fields{
			"method":     r.Method,
			"path":       r.URL.Path,
			"status":     lw.statusCode,
			"duration":   duration.String(),
			"user_agent": r.UserAgent(),
			"remote_ip":  r.RemoteAddr,
		}).Info("Request handled")
	})
}

// loggingResponseWriter is a wrapper around http.ResponseWriter that captures the status code.
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
EOF

# Update cmd/main.go
cat << EOF > cmd/main.go
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
EOF

# Run go mod tidy
go mod tidy

echo "Project setup complete in $PROJECT_DIR"
