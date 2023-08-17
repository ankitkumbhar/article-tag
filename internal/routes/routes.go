package routes

import (
	"article-tag/internal/handler"

	"github.com/go-chi/chi"
)

func InitRouter(app *handler.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/tag", app.Store())
	r.Get("/tag/{username}/{publication}", app.Get())
	r.Delete("/tag", app.Delete())
	r.Get("/tag/popular/{username}/{publication}", app.PopularTag())

	return r
}
