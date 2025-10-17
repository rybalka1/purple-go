#!/bin/bash

PROJECT_DIR="/home/rybalka/Myprojects/purple-go/5-product-api"

cd "$PROJECT_DIR/cmd"

# Run migrations (main.go handles this)
go run main.go
