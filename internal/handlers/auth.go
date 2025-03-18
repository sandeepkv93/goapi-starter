package handlers

import (
	"encoding/json"
	"goapi-starter/internal/config"
	"goapi-starter/internal/database"
	"goapi-starter/internal/models"
	"goapi-starter/internal/services"
	"goapi-starter/internal/utils"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	var req models.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Check if email already exists
	var existingUser models.User
	if result := database.DB.Where("email = ?", req.Email).First(&existingUser); result.Error == nil {
		utils.RespondWithError(w, http.StatusConflict, "Email already exists")
		return
	}

	// Check if username already exists
	if result := database.DB.Where("username = ?", req.Username).First(&existingUser); result.Error == nil {
		utils.RespondWithError(w, http.StatusConflict, "Username already exists")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
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
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

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
	var req models.SigninRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Find user by email
	var user models.User
	if result := database.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate token pair
	tokens, err := services.GenerateTokenPair(user)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error generating tokens")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, utils.SuccessResponse{
		Message: "Successfully signed in",
		Data: map[string]interface{}{
			"user":   user,
			"tokens": tokens,
		},
	})
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate refresh token and get user
	user, err := services.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	// Generate new token pair
	tokens, err := services.GenerateTokenPair(*user)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error generating tokens")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, utils.SuccessResponse{
		Message: "Tokens refreshed successfully",
		Data:    tokens,
	})
}

func generateJWT(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Second * time.Duration(config.AppConfig.JWT.AccessExpiry)).Unix(),
	})

	return token.SignedString([]byte(config.AppConfig.JWT.AccessSecret))
}
