package routes

import (
	"goapi-starter/internal/handlers"
	"goapi-starter/internal/utils"

	"github.com/go-chi/chi/v5"
)

func DummyProductRoutes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", utils.InstrumentHandler("CreateDummyProduct", handlers.CreateDummyProduct))
	r.Get("/", utils.InstrumentHandler("GetDummyProducts", handlers.GetDummyProducts))
	r.Get("/{id}", utils.InstrumentHandler("GetDummyProduct", handlers.GetDummyProduct))
	r.Put("/{id}", utils.InstrumentHandler("UpdateDummyProduct", handlers.UpdateDummyProduct))
	r.Delete("/{id}", utils.InstrumentHandler("DeleteDummyProduct", handlers.DeleteDummyProduct))

	return r
}
