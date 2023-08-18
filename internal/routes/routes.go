package routes

import (
	"article-tag/internal/handler"

	"github.com/go-chi/chi"
)

func InitRouter(app *handler.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/tags/{publication}", app.Store())
	r.Get("/tags/{publication}", app.Get())
	r.Delete("/tags/{publication}", app.Delete())
	r.Get("/tags/popular/{publication}", app.PopularTag())

	return r
}
