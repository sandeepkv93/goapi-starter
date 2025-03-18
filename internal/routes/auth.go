package routes

import (
	"goapi-starter/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes() chi.Router {
	r := chi.NewRouter()

	r.Post("/signup", handlers.SignUp)
	r.Post("/signin", handlers.SignIn)
	r.Post("/refresh", handlers.RefreshToken)

	return r
}
