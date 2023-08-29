package handler

import (
	"article-tag/internal/constant"
	"article-tag/internal/response"
	"article-tag/internal/types"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (app *Application) Store() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req types.StoreTagRequest

		// validate request
		err := app.validateStoreRequest(w, r, &req)
		if err != nil {
			app.logger.Error("error validating storing request", zap.Error(err), zap.Field{Key: "request",
				Type: zapcore.ReflectType, Interface: req})

			return
		}

		// store follow tag
		for _, val := range req.Tags {
			err = app.model.Tag.Store(ctx, req.Username, req.Publication, val.TagName, val.TagID)
			if err != nil {
				app.logger.Error("error while storing item", zap.Error(err), zap.Field{Key: "request",
					Type: zapcore.ReflectType, Interface: req})
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
			app.logger.Error("error validating get request", zap.Error(err), zap.Field{Key: "request",
				Type: zapcore.ReflectType, Interface: req})

			return
		}

		// fetch tags using username and publication
		userTags, err := app.model.Tag.Get(ctx, req.Username, req.Publication, req.Order)
		if err != nil {
			app.logger.Error("error fetching user tags from db", zap.Error(err), zap.Field{Key: "request",
				Type: zapcore.ReflectType, Interface: req})
			response.InternalServerError(w, "error while fetching user tags")

			return
		}

		var tags []types.Tag
		for _, val := range userTags {
			tags = append(tags, types.Tag{
				TagID:   val.TagID,
				TagName: val.TagName,
			})
		}

		// prepare response
		resp := types.GetTagResponse{Tags: tags}

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
			app.logger.Error("error validating delete request", zap.Error(err), zap.Field{Key: "request",
				Type: zapcore.ReflectType, Interface: req})

			return
		}

		for _, val := range req.Tags {
			// delete user tag
			err = app.model.Tag.Delete(ctx, req.Username, req.Publication, val.TagID, val.TagName)
			if err != nil {
				app.logger.Error("error deleting user tags", zap.Error(err), zap.Field{Key: "request",
					Type: zapcore.ReflectType, Interface: req})
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
			app.logger.Error("error validating get popular tag request", zap.Error(err), zap.Field{Key: "request",
				Type: zapcore.ReflectType, Interface: req})

			return
		}

		// fetch popularTags
		userTags, err := app.model.Tag.GetPopularTags(ctx, req.Username, req.Publication)
		if err != nil {
			app.logger.Error("error fetching popular tags from db", zap.Error(err), zap.Field{Key: "request",
				Type: zapcore.ReflectType, Interface: req})
			response.InternalServerError(w, "error while fetching popular tags")

			return
		}

		response.Success(w, userTags, "")
	}
}

func (app *Application) validateStoreRequest(w http.ResponseWriter, r *http.Request, req *types.StoreTagRequest) error {
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		app.logger.Error("error decoding store request body", zap.Error(err))
		response.BadRequest(w, "invalid request", nil)

		return err
	}

	req.Publication = chi.URLParam(r, "publication")

	// validator.InvalidValidationError

	err = app.validate.Struct(req)
	if err != nil {

		var errorBag []map[string]interface{}
		for _, v := range err.(validator.ValidationErrors) {
			errorBag = append(errorBag, map[string]interface{}{
				v.Field(): constant.TagError[v.Field()],
			})
		}

		response.BadRequest(w, "", errorBag)

		return err
	}

	return nil
}

func (app *Application) validateGetRequest(w http.ResponseWriter, r *http.Request, req *types.GetTagRequest) error {
	var err error

	// fetch username from queryParams
	req.Username = r.URL.Query().Get("username")

	req.Order = r.URL.Query().Get("order")

	// fetch params from urlParams
	req.Publication = chi.URLParam(r, "publication")

	err = app.validate.Struct(req)
	if err != nil {

		var errorBag []map[string]interface{}
		for _, v := range err.(validator.ValidationErrors) {
			errorBag = append(errorBag, map[string]interface{}{
				v.Field(): constant.TagError[v.Field()],
			})
		}

		response.BadRequest(w, "", errorBag)

		return err
	}

	return nil
}

func (app *Application) validateDeleteRequest(w http.ResponseWriter, r *http.Request, req *types.DeleteTagRequest) error {
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		app.logger.Error("error decoding delete request body", zap.Error(err))
		response.BadRequest(w, "invalid request", nil)

		return err
	}

	// fetch params from urlParams
	req.Publication = chi.URLParam(r, "publication")

	err = app.validate.Struct(req)
	if err != nil {

		var errorBag []map[string]interface{}
		for _, v := range err.(validator.ValidationErrors) {
			errorBag = append(errorBag, map[string]interface{}{
				v.Field(): constant.TagError[v.Field()],
			})
		}

		response.BadRequest(w, "", errorBag)

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

	err = app.validate.Struct(req)
	if err != nil {

		var errorBag []map[string]interface{}
		for _, v := range err.(validator.ValidationErrors) {
			errorBag = append(errorBag, map[string]interface{}{
				v.Field(): constant.TagError[v.Field()],
			})
		}

		response.BadRequest(w, "", errorBag)

		return err
	}

	return nil
}
