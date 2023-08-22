package handler

import (
	"article-tag/internal/constant"
	"article-tag/internal/model"
	"article-tag/internal/response"
	"article-tag/internal/types"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func (app *Application) Store() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req types.StoreTagRequest

		// validate request
		err := app.validateStoreRequest(w, r, &req)
		if err != nil {
			log.Println("error storing item : ", err)

			return
		}

		// check table exists
		err = app.model.Tag.DescribeTable(ctx)
		if err != nil {
			err = app.model.Tag.CreateTable(ctx)
			if err != nil {
				fmt.Println("error creating table : ", err)
				response.InternalServerError(w, "error while creating table")

				return
			}
		}

		// store follow tag
		for _, val := range req.Tags {
			item := model.UserTag{
				PK:          fmt.Sprintf("%v#%v", req.Username, req.Publication),
				SK:          val,
				Username:    req.Username,
				Publication: req.Publication,
				Tag:         val,
			}

			err = app.model.Tag.Store(ctx, &item)
			if err != nil {
				log.Println("error storing item : ", err)
				response.InternalServerError(w, "error while storing user tag")

				return
			}
		}

		response.Created(w, "")
	}
}

func (app *Application) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req types.GetTagRequest

		// validate request
		err := app.validateGetRequest(w, r, &req)
		if err != nil {
			log.Println("error storing item : ", err)

			return
		}

		item := model.UserTag{
			Username:    req.Username,
			Publication: req.Publication,
		}

		// fetch tags using username and publication
		userTag, err := app.model.Tag.Get(ctx, &item)
		if err != nil {
			log.Println("error fetching item : ", err)
			response.InternalServerError(w, "error while fetching user tags")

			return
		}

		// prepare response
		resp := types.GetTagResponse{Tags: userTag}

		response.Success(w, resp, "")
	}
}

func (app *Application) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req types.DeleteTagRequest

		// validate request
		err := app.validateDeleteRequest(w, r, &req)
		if err != nil {
			log.Println("error storing item : ", err)

			return
		}

		for _, val := range req.Tags {
			item := model.UserTag{
				Username:    req.Username,
				Publication: req.Publication,
				Tag:         val,
			}

			err = app.model.Tag.Delete(ctx, &item)
			if err != nil {
				log.Println("error deleting item : ", err)
				response.InternalServerError(w, "error while deleting user followed tags")

				return
			}
		}

		response.Success(w, nil, "")
	}
}

func (app *Application) PopularTag() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req types.GetPopularTagRequest

		// validate request
		err := app.validateGetPopularTagRequest(w, r, &req)
		if err != nil {
			log.Println("error storing item : ", err)

			return
		}

		item := model.UserTag{
			Username:    req.Username,
			Publication: req.Publication,
		}

		userTags, err := app.model.Tag.GetPopularTags(ctx, &item)
		if err != nil {
			log.Println("error while fetching popular tags for user", err)
			response.InternalServerError(w, "error while fetching popular tags")

			return
		}

		response.Success(w, userTags, "")
	}
}

func (app *Application) validateStoreRequest(w http.ResponseWriter, r *http.Request, req *types.StoreTagRequest) error {
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("error decoding request body : ", err)
		response.BadRequest(w, "invalid request")

		return err
	}

	req.Publication = chi.URLParam(r, "publication")

	// TODO : use validator to validate request

	if req.Username == "" {
		err = errors.New("username field is required")
		response.BadRequest(w, "username field is required")

		return err
	}

	if req.Publication == "" {
		err = errors.New("publication field is required")
		response.BadRequest(w, "publication field is required")

		return err
	}

	var isValidPublication bool
	for _, v := range constant.AllowdedPublications {
		if req.Publication == v {
			isValidPublication = true

			break
		}
	}

	if !isValidPublication {
		err = errors.New("invalid publication passed")
		response.BadRequest(w, "invalid publication passed, please pass valid publication")

		return err
	}

	if len(req.Tags) == 0 {
		err = errors.New("atleast one tag is required")
		response.BadRequest(w, "atleast one tag is required")

		return err
	}

	return nil
}

func (app *Application) validateGetRequest(w http.ResponseWriter, r *http.Request, req *types.GetTagRequest) error {
	var err error

	// fetch username from queryParams
	req.Username = r.URL.Query().Get("username")

	// fetch params from urlParams
	req.Publication = chi.URLParam(r, "publication")

	// TODO : use validator to validate request

	if req.Username == "" {
		err = errors.New("username field is required")
		response.BadRequest(w, "username field is required")

		return err
	}

	if req.Publication == "" {
		err = errors.New("publication field is required")
		response.BadRequest(w, "publication field is required")

		return err
	}

	return nil
}

func (app *Application) validateDeleteRequest(w http.ResponseWriter, r *http.Request, req *types.DeleteTagRequest) error {
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("error decoding request body : ", err)
		response.BadRequest(w, "invalid request")

		return err
	}

	// fetch params from urlParams
	req.Publication = chi.URLParam(r, "publication")

	// TODO : use validator to validate request

	if req.Username == "" {
		err = errors.New("username field is required")
		response.BadRequest(w, "username field is required")

		return err
	}

	if req.Publication == "" {
		err = errors.New("publication field is required")
		response.BadRequest(w, "publication field is required")

		return err
	}

	var isValidPublication bool
	for _, v := range constant.AllowdedPublications {
		if req.Publication == v {
			isValidPublication = true

			break
		}
	}

	if !isValidPublication {
		err = errors.New("invalid publication passed")
		response.BadRequest(w, "invalid publication passed, please pass valid publication")

		return err
	}

	if len(req.Tags) == 0 {
		err = errors.New("atleast one tag is required")
		response.BadRequest(w, "atleast one tag is required")

		return err
	}

	return nil
}

func (app *Application) validateGetPopularTagRequest(w http.ResponseWriter, r *http.Request, req *types.GetPopularTagRequest) error {
	var err error

	// fetch username from queryParams
	req.Username = r.URL.Query().Get("username")

	// fetch params from urlParams
	req.Publication = chi.URLParam(r, "publication")

	// TODO : use validator to validate request

	if req.Publication == "" {
		err = errors.New("publication field is required")
		response.BadRequest(w, "publication field is required")

		return err
	}

	return nil
}
