package routes

import (
	"article-tag/internal/handler"
	"article-tag/internal/response"
	"net/http"

	"github.com/go-chi/chi"
)

func InitRouter(app *handler.Application) *chi.Mux {
	r := chi.NewRouter()

	// sets a custom message for 404 error status code
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		response.NotFound(w, "requested url is unavailable")
		return
	})

	// sets custom message for 405 error status code
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		response.NotAllowded(w, "method not allowed")
		return
	})

	// route group
	r.Route("/tags", func(r chi.Router) {
		r.Post("/{publication}", app.Store())
		r.Get("/{publication}", app.Get())
		r.Delete("/{publication}", app.Delete())
		r.Get("/{publication}/popular", app.PopularTag())
	})

	return r
}
