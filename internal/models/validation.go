package models

type ProductRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=100"`
	Description string  `json:"description" validate:"max=500"`
	Price       float64 `json:"price" validate:"required,gt=0"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name" validate:"omitempty,min=3,max=100"`
	Description *string  `json:"description" validate:"omitempty,max=500"`
	Price       *float64 `json:"price" validate:"omitempty,gt=0"`
}

type SignupRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type SigninRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // seconds until access token expires
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
