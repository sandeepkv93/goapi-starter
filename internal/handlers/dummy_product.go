package handlers

import (
	"encoding/json"
	"fmt"
	"goapi-starter/internal/cache"
	"goapi-starter/internal/database"
	"goapi-starter/internal/logger"
	"goapi-starter/internal/metrics"
	"goapi-starter/internal/models"
	"goapi-starter/internal/utils"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// CreateDummyProduct handles the creation of a new dummy product
func CreateDummyProduct(w http.ResponseWriter, r *http.Request) {
	metrics.BusinessOperations.WithLabelValues("create_dummy_product", "started").Inc()

	var req models.DummyProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		metrics.RecordHandlerError("CreateDummyProduct", "invalid_request")
		metrics.RecordDetailedError("CreateDummyProduct", "invalid_request", "json_decode_error")
		metrics.BusinessOperations.WithLabelValues("create_dummy_product", "failed").Inc()
		utils.RespondWithError(w, r, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		metrics.RecordHandlerError("CreateDummyProduct", "validation_error")
		metrics.RecordDetailedError("CreateDummyProduct", "validation_error", err.Error())
		metrics.BusinessOperations.WithLabelValues("create_dummy_product", "failed").Inc()
		utils.RespondWithError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	dummyProduct := models.DummyProduct{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}

	result := database.DB.Create(&dummyProduct)
	if result.Error != nil {
		errorReason := "unknown"
		if result.Error != nil {
			if strings.Contains(result.Error.Error(), "duplicate") {
				errorReason = "duplicate_entry"
			} else {
				// Limit the error reason length to avoid cardinality explosion
				if len(result.Error.Error()) > 50 {
					errorReason = result.Error.Error()[:50]
				} else {
					errorReason = result.Error.Error()
				}
			}
		}

		metrics.RecordHandlerError("CreateDummyProduct", "database_error")
		metrics.RecordDetailedError("CreateDummyProduct", "database_error", errorReason)
		metrics.BusinessOperations.WithLabelValues("create_dummy_product", "failed").Inc()
		utils.RespondWithError(w, r, http.StatusInternalServerError, "Error creating dummy product")
		return
	}

	// Invalidate the list cache since we've added a new product
	if err := cache.Delete("dummy_products:all"); err != nil {
		logger.Warn().Err(err).Msg("Failed to invalidate dummy products list cache")
	}

	metrics.BusinessOperations.WithLabelValues("create_dummy_product", "success").Inc()
	utils.RespondWithJSON(w, r, http.StatusCreated, utils.SuccessResponse{
		Message: "Dummy product created successfully",
		Data:    dummyProduct,
	})
}

// GetDummyProducts returns a list of all dummy products
func GetDummyProducts(w http.ResponseWriter, r *http.Request) {
	metrics.BusinessOperations.WithLabelValues("get_dummy_products", "started").Inc()

	// Try to get from cache first
	var dummyProducts []models.DummyProduct
	cacheKey := "dummy_products:all"

	found, err := cache.Get(cacheKey, &dummyProducts)
	if err != nil {
		logger.Warn().Err(err).Msg("Error retrieving from cache")
		// Continue with database query
	}

	if found && len(dummyProducts) > 0 {
		logger.Info().Msg("Returning dummy products from cache")
		metrics.BusinessOperations.WithLabelValues("get_dummy_products", "success").Inc()
		utils.RespondWithJSON(w, r, http.StatusOK, utils.SuccessResponse{
			Message: "Dummy products retrieved from cache",
			Data:    dummyProducts,
		})
		return
	}

	// Not in cache, get from database
	result := database.DB.Find(&dummyProducts)
	if result.Error != nil {
		metrics.RecordHandlerError("GetDummyProducts", "database_error")
		metrics.RecordDetailedError("GetDummyProducts", "database_error", result.Error.Error())
		metrics.BusinessOperations.WithLabelValues("get_dummy_products", "failed").Inc()
		utils.RespondWithError(w, r, http.StatusInternalServerError, "Error retrieving dummy products")
		return
	}

	// Store in cache for future requests
	if len(dummyProducts) > 0 {
		if err := cache.Set(cacheKey, dummyProducts); err != nil {
			logger.Warn().Err(err).Msg("Failed to cache dummy products")
		}
	}

	metrics.BusinessOperations.WithLabelValues("get_dummy_products", "success").Inc()
	utils.RespondWithJSON(w, r, http.StatusOK, utils.SuccessResponse{
		Message: "Dummy products retrieved successfully",
		Data:    dummyProducts,
	})
}

// GetDummyProduct returns a specific dummy product by ID
func GetDummyProduct(w http.ResponseWriter, r *http.Request) {
	metrics.BusinessOperations.WithLabelValues("get_dummy_product", "started").Inc()

	id := chi.URLParam(r, "id")
	if id == "" {
		metrics.RecordHandlerError("GetDummyProduct", "invalid_request")
		metrics.RecordDetailedError("GetDummyProduct", "invalid_request", "missing_id")
		metrics.BusinessOperations.WithLabelValues("get_dummy_product", "failed").Inc()
		utils.RespondWithError(w, r, http.StatusBadRequest, "Missing dummy product ID")
		return
	}

	// Try to get from cache first
	var dummyProduct models.DummyProduct
	cacheKey := fmt.Sprintf("dummy_product:%s", id)

	found, err := cache.Get(cacheKey, &dummyProduct)
	if err != nil {
		logger.Warn().Err(err).Str("id", id).Msg("Error retrieving product from cache")
		// Continue with database query
	}

	if found {
		logger.Info().Str("id", id).Msg("Returning dummy product from cache")
		metrics.BusinessOperations.WithLabelValues("get_dummy_product", "success").Inc()
		utils.RespondWithJSON(w, r, http.StatusOK, utils.SuccessResponse{
			Message: "Dummy product retrieved from cache",
			Data:    dummyProduct,
		})
		return
	}

	// Not in cache, get from database
	result := database.DB.First(&dummyProduct, id)
	if result.Error != nil {
		metrics.RecordHandlerError("GetDummyProduct", "not_found")
		metrics.RecordDetailedError("GetDummyProduct", "not_found", "id_"+id)
		metrics.BusinessOperations.WithLabelValues("get_dummy_product", "failed").Inc()
		utils.RespondWithError(w, r, http.StatusNotFound, "Dummy product not found")
		return
	}

	// Store in cache for future requests
	if err := cache.Set(cacheKey, dummyProduct); err != nil {
		logger.Warn().Err(err).Str("id", id).Msg("Failed to cache dummy product")
	}

	metrics.BusinessOperations.WithLabelValues("get_dummy_product", "success").Inc()
	utils.RespondWithJSON(w, r, http.StatusOK, utils.SuccessResponse{
		Message: "Dummy product retrieved successfully",
		Data:    dummyProduct,
	})
}

// UpdateDummyProduct updates a specific dummy product
func UpdateDummyProduct(w http.ResponseWriter, r *http.Request) {
	metrics.BusinessOperations.WithLabelValues("update_dummy_product", "started").Inc()

	id := chi.URLParam(r, "id")
	if id == "" {
		metrics.RecordHandlerError("UpdateDummyProduct", "invalid_request")
		metrics.RecordDetailedError("UpdateDummyProduct", "invalid_request", "missing_id")
		metrics.BusinessOperations.WithLabelValues("update_dummy_product", "failed").Inc()
		utils.RespondWithError(w, r, http.StatusBadRequest, "Missing dummy product ID")
		return
	}

	var req models.UpdateDummyProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		metrics.RecordHandlerError("UpdateDummyProduct", "invalid_request")
		metrics.RecordDetailedError("UpdateDummyProduct", "invalid_request", "json_decode_error")
		metrics.BusinessOperations.WithLabelValues("update_dummy_product", "failed").Inc()
		utils.RespondWithError(w, r, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		metrics.RecordHandlerError("UpdateDummyProduct", "validation_error")
		metrics.RecordDetailedError("UpdateDummyProduct", "validation_error", err.Error())
		metrics.BusinessOperations.WithLabelValues("update_dummy_product", "failed").Inc()
		utils.RespondWithError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// Check if dummy product exists
	var dummyProduct models.DummyProduct
	if result := database.DB.First(&dummyProduct, id); result.Error != nil {
		metrics.RecordHandlerError("UpdateDummyProduct", "not_found")
		metrics.RecordDetailedError("UpdateDummyProduct", "not_found", "id_"+id)
		metrics.BusinessOperations.WithLabelValues("update_dummy_product", "failed").Inc()
		utils.RespondWithError(w, r, http.StatusNotFound, "Dummy product not found")
		return
	}

	// Update fields if provided
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}

	if len(updates) == 0 {
		metrics.RecordHandlerError("UpdateDummyProduct", "invalid_request")
		metrics.RecordDetailedError("UpdateDummyProduct", "invalid_request", "no_updates")
		metrics.BusinessOperations.WithLabelValues("update_dummy_product", "failed").Inc()
		utils.RespondWithError(w, r, http.StatusBadRequest, "No updates provided")
		return
	}

	if result := database.DB.Model(&dummyProduct).Updates(updates); result.Error != nil {
		metrics.RecordHandlerError("UpdateDummyProduct", "database_error")
		metrics.RecordDetailedError("UpdateDummyProduct", "database_error", result.Error.Error())
		metrics.BusinessOperations.WithLabelValues("update_dummy_product", "failed").Inc()
		utils.RespondWithError(w, r, http.StatusInternalServerError, "Error updating dummy product")
		return
	}

	// Get the updated dummy product
	database.DB.First(&dummyProduct, id)

	// Update the product in cache
	cacheKey := fmt.Sprintf("dummy_product:%s", id)
	if err := cache.Set(cacheKey, dummyProduct); err != nil {
		logger.Warn().Err(err).Str("id", id).Msg("Failed to update dummy product in cache")
	}

	// Invalidate the list cache since a product was updated
	if err := cache.Delete("dummy_products:all"); err != nil {
		logger.Warn().Err(err).Msg("Failed to invalidate dummy products list cache")
	}

	metrics.BusinessOperations.WithLabelValues("update_dummy_product", "success").Inc()
	utils.RespondWithJSON(w, r, http.StatusOK, utils.SuccessResponse{
		Message: "Dummy product updated successfully",
		Data:    dummyProduct,
	})
}

// DeleteDummyProduct deletes a specific dummy product
func DeleteDummyProduct(w http.ResponseWriter, r *http.Request) {
	metrics.BusinessOperations.WithLabelValues("delete_dummy_product", "started").Inc()

	id := chi.URLParam(r, "id")
	if id == "" {
		metrics.RecordHandlerError("DeleteDummyProduct", "invalid_request")
		metrics.RecordDetailedError("DeleteDummyProduct", "invalid_request", "missing_id")
		metrics.BusinessOperations.WithLabelValues("delete_dummy_product", "failed").Inc()
		utils.RespondWithError(w, r, http.StatusBadRequest, "Missing dummy product ID")
		return
	}

	// Check if dummy product exists
	var dummyProduct models.DummyProduct
	if result := database.DB.First(&dummyProduct, id); result.Error != nil {
		metrics.RecordHandlerError("DeleteDummyProduct", "not_found")
		metrics.RecordDetailedError("DeleteDummyProduct", "not_found", "id_"+id)
		metrics.BusinessOperations.WithLabelValues("delete_dummy_product", "failed").Inc()
		utils.RespondWithError(w, r, http.StatusNotFound, "Dummy product not found")
		return
	}

	// Delete the dummy product
	if result := database.DB.Delete(&dummyProduct); result.Error != nil {
		metrics.RecordHandlerError("DeleteDummyProduct", "database_error")
		metrics.RecordDetailedError("DeleteDummyProduct", "database_error", result.Error.Error())
		metrics.BusinessOperations.WithLabelValues("delete_dummy_product", "failed").Inc()
		utils.RespondWithError(w, r, http.StatusInternalServerError, "Error deleting dummy product")
		return
	}

	// Delete the product from cache
	cacheKey := fmt.Sprintf("dummy_product:%s", id)
	if err := cache.Delete(cacheKey); err != nil {
		logger.Warn().Err(err).Str("id", id).Msg("Failed to delete dummy product from cache")
	}

	// Invalidate the list cache since a product was deleted
	if err := cache.Delete("dummy_products:all"); err != nil {
		logger.Warn().Err(err).Msg("Failed to invalidate dummy products list cache")
	}

	metrics.BusinessOperations.WithLabelValues("delete_dummy_product", "success").Inc()
	utils.RespondWithJSON(w, r, http.StatusOK, utils.SuccessResponse{
		Message: "Dummy product deleted successfully",
		Data:    nil,
	})
}
