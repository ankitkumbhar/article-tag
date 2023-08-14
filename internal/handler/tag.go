package handler

import (
	"article-tag/internal/response"
	"fmt"
	"net/http"
)

func (app *Application) Store() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := app.model.Tag.Store(ctx, nil)
		if err != nil {
			fmt.Println("error storing item : ", err)
			response.InternalServerError(w, "error while storing item from database")

			return
		}

		response.Created(w, "")
	}
}
