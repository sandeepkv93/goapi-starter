package routes

import (
	"goapi-starter/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func UserRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", handlers.ListUsers)
	r.Post("/", handlers.CreateUser)
	r.Get("/{id}", handlers.GetUser)
	// Add more routes as needed

	return r
}
