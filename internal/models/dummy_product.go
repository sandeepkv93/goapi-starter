package models

import (
	"time"
)

// DummyProduct represents a dummy product in the system
type DummyProduct struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Description string    `json:"description" gorm:"size:500"`
	Price       float64   `json:"price" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// DummyProductRequest is used for creating or updating a dummy product
type DummyProductRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=100"`
	Description string  `json:"description" validate:"max=500"`
	Price       float64 `json:"price" validate:"required,gt=0"`
}

// UpdateDummyProductRequest is used for partial updates to a dummy product
type UpdateDummyProductRequest struct {
	Name        *string  `json:"name" validate:"omitempty,min=3,max=100"`
	Description *string  `json:"description" validate:"omitempty,max=500"`
	Price       *float64 `json:"price" validate:"omitempty,gt=0"`
}
