package routes

import (
	"cursor-experiment-1/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func ProductRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", handlers.ListProducts)
	r.Post("/", handlers.CreateProduct)
	r.Get("/{id}", handlers.GetProduct)
	r.Put("/{id}", handlers.UpdateProduct)
	r.Delete("/{id}", handlers.DeleteProduct)

	return r
}
