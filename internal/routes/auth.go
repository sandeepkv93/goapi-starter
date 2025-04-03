package routes

import (
	"goapi-starter/internal/handlers"
	"goapi-starter/internal/utils"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes() chi.Router {
	r := chi.NewRouter()

	r.Post("/signup", utils.InstrumentHandler("SignUp", handlers.SignUp))
	r.Post("/signin", utils.InstrumentHandler("SignIn", handlers.SignIn))
	r.Post("/refresh", utils.InstrumentHandler("RefreshToken", handlers.RefreshToken))

	return r
}
