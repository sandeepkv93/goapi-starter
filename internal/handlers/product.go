package handlers

import (
	"encoding/json"
	"goapi-starter/internal/database"
	"goapi-starter/internal/metrics"
	"goapi-starter/internal/models"
	"goapi-starter/internal/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func ListProducts(w http.ResponseWriter, r *http.Request) {
	metrics.BusinessOperations.WithLabelValues("list_products", "started").Inc()

	var products []models.Product
	result := database.DB.Find(&products)
	if result.Error != nil {
		metrics.RecordHandlerError("ListProducts", "database_error")
		metrics.BusinessOperations.WithLabelValues("list_products", "failed").Inc()
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching products")
		return
	}

	metrics.BusinessOperations.WithLabelValues("list_products", "success").Inc()
	utils.RespondWithJSON(w, http.StatusOK, utils.SuccessResponse{
		Data: products,
	})
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	metrics.BusinessOperations.WithLabelValues("create_product", "started").Inc()

	var req models.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		metrics.RecordHandlerError("CreateProduct", "invalid_request")
		metrics.BusinessOperations.WithLabelValues("create_product", "failed").Inc()
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		metrics.RecordHandlerError("CreateProduct", "validation_error")
		metrics.BusinessOperations.WithLabelValues("create_product", "failed").Inc()
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
		metrics.RecordHandlerError("CreateProduct", "database_error")
		metrics.BusinessOperations.WithLabelValues("create_product", "failed").Inc()
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating product")
		return
	}

	metrics.BusinessOperations.WithLabelValues("create_product", "success").Inc()
	utils.RespondWithJSON(w, http.StatusCreated, utils.SuccessResponse{
		Message: "Product created successfully",
		Data:    product,
	})
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	metrics.BusinessOperations.WithLabelValues("get_product", "started").Inc()

	id := chi.URLParam(r, "id")

	var product models.Product
	result := database.DB.First(&product, "id = ?", id)
	if result.Error != nil {
		metrics.RecordHandlerError("GetProduct", "not_found")
		metrics.BusinessOperations.WithLabelValues("get_product", "failed").Inc()
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	metrics.BusinessOperations.WithLabelValues("get_product", "success").Inc()
	utils.RespondWithJSON(w, http.StatusOK, utils.SuccessResponse{
		Data: product,
	})
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	metrics.BusinessOperations.WithLabelValues("update_product", "started").Inc()

	id := chi.URLParam(r, "id")

	var product models.Product
	result := database.DB.First(&product, "id = ?", id)
	if result.Error != nil {
		metrics.RecordHandlerError("UpdateProduct", "not_found")
		metrics.BusinessOperations.WithLabelValues("update_product", "failed").Inc()
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	var req models.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		metrics.RecordHandlerError("UpdateProduct", "invalid_request")
		metrics.BusinessOperations.WithLabelValues("update_product", "failed").Inc()
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		metrics.RecordHandlerError("UpdateProduct", "validation_error")
		metrics.BusinessOperations.WithLabelValues("update_product", "failed").Inc()
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price

	result = database.DB.Save(&product)
	if result.Error != nil {
		metrics.RecordHandlerError("UpdateProduct", "database_error")
		metrics.BusinessOperations.WithLabelValues("update_product", "failed").Inc()
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating product")
		return
	}

	metrics.BusinessOperations.WithLabelValues("update_product", "success").Inc()
	utils.RespondWithJSON(w, http.StatusOK, utils.SuccessResponse{
		Message: "Product updated successfully",
		Data:    product,
	})
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	metrics.BusinessOperations.WithLabelValues("delete_product", "started").Inc()

	id := chi.URLParam(r, "id")

	var product models.Product
	result := database.DB.First(&product, "id = ?", id)
	if result.Error != nil {
		metrics.RecordHandlerError("DeleteProduct", "not_found")
		metrics.BusinessOperations.WithLabelValues("delete_product", "failed").Inc()
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	result = database.DB.Delete(&product)
	if result.Error != nil {
		metrics.RecordHandlerError("DeleteProduct", "database_error")
		metrics.BusinessOperations.WithLabelValues("delete_product", "failed").Inc()
		utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting product")
		return
	}

	metrics.BusinessOperations.WithLabelValues("delete_product", "success").Inc()
	utils.RespondWithJSON(w, http.StatusOK, utils.SuccessResponse{
		Message: "Product deleted successfully",
	})
}
