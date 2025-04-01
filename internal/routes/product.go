package routes

import (
	"goapi-starter/internal/handlers"
	"goapi-starter/internal/utils"

	"github.com/go-chi/chi/v5"
)

func ProductRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", utils.InstrumentHandler("ListProducts", handlers.ListProducts))
	r.Post("/", utils.InstrumentHandler("CreateProduct", handlers.CreateProduct))
	r.Get("/{id}", utils.InstrumentHandler("GetProduct", handlers.GetProduct))
	r.Put("/{id}", utils.InstrumentHandler("UpdateProduct", handlers.UpdateProduct))
	r.Delete("/{id}", utils.InstrumentHandler("DeleteProduct", handlers.DeleteProduct))

	return r
}
