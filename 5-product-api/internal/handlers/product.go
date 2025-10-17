package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"purple-go/5-product-api/internal/database"
	"purple-go/5-product-api/internal/models"
	"purple-go/5-product-api/internal/utils"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := database.DB.Create(&product).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to create product")
		return
	}

	utils.JSON(w, http.StatusCreated, product)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var product models.Product
	if err := database.DB.First(&product, params["id"]).Error; err != nil {
		utils.Error(w, http.StatusNotFound, "Product not found")
		return
	}
	utils.JSON(w, http.StatusOK, product)
}

func UpdateProduct(w http.ResponseWriter, r *request) {
	params := mux.Vars(r)
	var product models.Product
	if err := database.DB.First(&product, params["id"]).Error; err != nil {
		utils.Error(w, http.StatusNotFound, "Product not found")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	database.DB.Save(&product)
	utils.JSON(w, http.StatusOK, product)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var product models.Product
	if err := database.DB.First(&product, params["id"]).Error; err != nil {
		utils.Error(w, http.StatusNotFound, "Product not found")
		return
	}

	database.DB.Delete(&product)
	utils.JSON(w, http.StatusOK, map[string]string{"result": "success"})
}
