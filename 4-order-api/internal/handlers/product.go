package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rybalka1/purple-go/4-order-api/internal/database"
	"github.com/rybalka1/purple-go/4-order-api/internal/models"
	"github.com/rybalka1/purple-go/4-order-api/internal/utils"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate required fields
	fields := map[string]string{
		"Name":        product.Name,
		"Description": product.Description,
	}
	if !utils.ValidateFields(fields, []string{"Name", "Description"}, w) {
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

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
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

	// Validate required fields
	fields := map[string]string{
		"Name":        product.Name,
		"Description": product.Description,
	}
	if !utils.ValidateFields(fields, []string{"Name", "Description"}, w) {
		return
	}

	if err := database.DB.Save(&product).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to update product")
		return
	}
	utils.JSON(w, http.StatusOK, product)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var product models.Product
	if err := database.DB.First(&product, params["id"]).Error; err != nil {
		utils.Error(w, http.StatusNotFound, "Product not found")
		return
	}

	if err := database.DB.Delete(&product).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to delete product")
		return
	}
	utils.JSON(w, http.StatusOK, map[string]string{"result": "success"})
}
