package routes

import (
	"article-tag/internal/handler"

	"github.com/go-chi/chi"
)

func InitRouter(app *handler.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/tag", app.Store())

	return r
}
