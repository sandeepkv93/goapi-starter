package routes

import (
	"goapi-starter/internal/handlers"
	"goapi-starter/internal/middleware"
	"goapi-starter/internal/utils"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes() chi.Router {
	r := chi.NewRouter()

	// Public auth endpoints (signup, signin)
	r.Group(func(r chi.Router) {
		r.Post("/signup", utils.InstrumentHandler("SignUp", handlers.SignUp))
		r.Post("/signin", utils.InstrumentHandler("SignIn", handlers.SignIn))
	})

	// Protected auth endpoints (refresh, logout)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Post("/refresh", utils.InstrumentHandler("RefreshToken", handlers.RefreshToken))
		r.Post("/logout", utils.InstrumentHandler("Logout", handlers.Logout))
	})

	return r
}
