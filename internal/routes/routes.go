package routes

import (
	"article-tag/internal/handler"
	"article-tag/internal/response"
	"net/http"

	"github.com/go-chi/chi"
)

func InitRouter(app *handler.Application) *chi.Mux {
	r := chi.NewRouter()

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		response.NotFound(w, "requested url is unavailable")
		return
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		response.NotAllowded(w, "method not allowed")
		return
	})

	r.Post("/tags/{publication}", app.Store())
	r.Get("/tags/{publication}", app.Get())
	r.Delete("/tags/{publication}", app.Delete())
	r.Get("/tags/{publication}/popular", app.PopularTag())

	return r
}
