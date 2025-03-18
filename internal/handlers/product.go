package handlers

import (
	"goapi-starter/internal/database"
	"goapi-starter/internal/models"
	"goapi-starter/internal/utils"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func ListProducts(w http.ResponseWriter, r *http.Request) {
	var products []models.Product
	result := database.DB.Find(&products)
	if result.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching products")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, utils.SuccessResponse{
		Data: products,
	})
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req models.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	product := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}

	result := database.DB.Create(&product)
	if result.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating product")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, utils.SuccessResponse{
		Message: "Product created successfully",
		Data:    product,
	})
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var product models.Product
	result := database.DB.First(&product, "id = ?", id)
	if result.Error != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, utils.SuccessResponse{
		Data: product,
	})
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var product models.Product
	result := database.DB.First(&product, "id = ?", id)
	if result.Error != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	// Update only provided fields
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}

	result = database.DB.Save(&product)
	if result.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating product")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, utils.SuccessResponse{
		Message: "Product updated successfully",
		Data:    product,
	})
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result := database.DB.Delete(&models.Product{}, "id = ?", id)
	if result.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting product")
		return
	}

	if result.RowsAffected == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, utils.SuccessResponse{
		Message: "Product deleted successfully",
	})
}
