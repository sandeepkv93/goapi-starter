package handlers

import (
	"encoding/json"
	"goapi-starter/internal/database"
	"goapi-starter/internal/metrics"
	"goapi-starter/internal/models"
	"goapi-starter/internal/services"
	"goapi-starter/internal/utils"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	metrics.BusinessOperations.WithLabelValues("signup", "started").Inc()

	var req models.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		metrics.RecordHandlerError("SignUp", "invalid_request")
		metrics.BusinessOperations.WithLabelValues("signup", "failed").Inc()
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		metrics.RecordHandlerError("SignUp", "validation_error")
		metrics.BusinessOperations.WithLabelValues("signup", "failed").Inc()
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Check if email already exists
	var existingUser models.User
	if result := database.DB.Where("email = ?", req.Email).First(&existingUser); result.Error == nil {
		metrics.RecordHandlerError("SignUp", "email_exists")
		metrics.BusinessOperations.WithLabelValues("signup", "failed").Inc()
		utils.RespondWithError(w, http.StatusConflict, "Email already exists")
		return
	}

	// Check if username already exists
	if result := database.DB.Where("username = ?", req.Username).First(&existingUser); result.Error == nil {
		metrics.RecordHandlerError("SignUp", "username_exists")
		metrics.BusinessOperations.WithLabelValues("signup", "failed").Inc()
		utils.RespondWithError(w, http.StatusConflict, "Username already exists")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		metrics.RecordHandlerError("SignUp", "password_hash_error")
		metrics.BusinessOperations.WithLabelValues("signup", "failed").Inc()
		utils.RespondWithError(w, http.StatusInternalServerError, "Error processing request")
		return
	}

	// Create user
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if result := database.DB.Create(&user); result.Error != nil {
		metrics.RecordHandlerError("SignUp", "database_error")
		metrics.BusinessOperations.WithLabelValues("signup", "failed").Inc()
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	metrics.BusinessOperations.WithLabelValues("signup", "success").Inc()
	// Return user data without password
	utils.RespondWithJSON(w, http.StatusCreated, utils.SuccessResponse{
		Message: "User created successfully",
		Data: map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	metrics.BusinessOperations.WithLabelValues("signin", "started").Inc()

	var req models.SigninRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		metrics.RecordHandlerError("SignIn", "invalid_request")
		metrics.BusinessOperations.WithLabelValues("signin", "failed").Inc()
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		metrics.RecordHandlerError("SignIn", "validation_error")
		metrics.BusinessOperations.WithLabelValues("signin", "failed").Inc()
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Find user by email
	var user models.User
	if result := database.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		metrics.RecordHandlerError("SignIn", "user_not_found")
		metrics.BusinessOperations.WithLabelValues("signin", "failed").Inc()
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		metrics.RecordHandlerError("SignIn", "invalid_password")
		metrics.BusinessOperations.WithLabelValues("signin", "failed").Inc()
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate token pair
	tokens, err := services.GenerateTokenPair(user)
	if err != nil {
		metrics.RecordHandlerError("SignIn", "token_generation_error")
		metrics.BusinessOperations.WithLabelValues("signin", "failed").Inc()
		utils.RespondWithError(w, http.StatusInternalServerError, "Error generating tokens")
		return
	}

	metrics.BusinessOperations.WithLabelValues("signin", "success").Inc()
	utils.RespondWithJSON(w, http.StatusOK, utils.SuccessResponse{
		Message: "Successfully signed in",
		Data: map[string]interface{}{
			"user":   user,
			"tokens": tokens,
		},
	})
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	metrics.BusinessOperations.WithLabelValues("refresh_token", "started").Inc()

	var req models.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		metrics.RecordHandlerError("RefreshToken", "invalid_request")
		metrics.BusinessOperations.WithLabelValues("refresh_token", "failed").Inc()
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		metrics.RecordHandlerError("RefreshToken", "validation_error")
		metrics.BusinessOperations.WithLabelValues("refresh_token", "failed").Inc()
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate refresh token and get user
	user, err := services.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		metrics.RecordHandlerError("RefreshToken", "invalid_token")
		metrics.BusinessOperations.WithLabelValues("refresh_token", "failed").Inc()
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	// Generate new token pair
	tokens, err := services.GenerateTokenPair(*user)
	if err != nil {
		metrics.RecordHandlerError("RefreshToken", "token_generation_error")
		metrics.BusinessOperations.WithLabelValues("refresh_token", "failed").Inc()
		utils.RespondWithError(w, http.StatusInternalServerError, "Error generating tokens")
		return
	}

	metrics.BusinessOperations.WithLabelValues("refresh_token", "success").Inc()
	utils.RespondWithJSON(w, http.StatusOK, utils.SuccessResponse{
		Message: "Tokens refreshed successfully",
		Data:    tokens,
	})
}
