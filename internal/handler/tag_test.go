package handler_test

import (
	"article-tag/internal/handler"
	"article-tag/internal/mocks"
	"article-tag/internal/model"
	"article-tag/internal/response"
	"article-tag/internal/types"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Store(t *testing.T) {
	type args struct {
		req types.StoreTagRequest
	}

	tests := []struct {
		name         string
		args         args
		mockDB       func() *handler.Application
		wantResp     string
		wantRespBody *response.Body
	}{
		{
			name: "success",
			args: args{req: types.StoreTagRequest{Username: "Ankit", Publication: "AK", Tags: []string{"tag100"}}},
			mockDB: func() *handler.Application {
				tagStoreMock := mocks.NewUserTagStore(t)
				tagStoreMock.EXPECT().Store(mock.Anything, mock.Anything).Return(nil)

				m := model.Models{
					Tag: tagStoreMock,
				}

				return handler.New(nil, &m)
			},
			// wantResp:     handler.ArticleResponse{ID: 1},
			wantRespBody: &response.Body{Status: http.StatusCreated, Message: ""},
		},
		{
			name: "should fail when invalid request is passed - empty username",
			args: args{req: types.StoreTagRequest{Username: "", Publication: "AK", Tags: []string{"tag100"}}},
			mockDB: func() *handler.Application {
				m := model.Models{}

				return handler.New(nil, &m)
			},
			wantRespBody: &response.Body{Status: http.StatusBadRequest, Message: "username field is required"},
		},
		{
			name: "should fail when invalid request is passed - empty publication",
			args: args{req: types.StoreTagRequest{Username: "Ankit", Publication: "", Tags: []string{"tag100"}}},
			mockDB: func() *handler.Application {
				m := model.Models{}

				return handler.New(nil, &m)
			},
			wantRespBody: &response.Body{Status: http.StatusBadRequest, Message: "publication field is required"},
		},
		{
			name: "should fail when invalid request is passed - empty tags",
			args: args{req: types.StoreTagRequest{Username: "Ankit", Publication: "AK", Tags: []string{}}},
			mockDB: func() *handler.Application {
				m := model.Models{}

				return handler.New(nil, &m)
			},
			wantRespBody: &response.Body{Status: http.StatusBadRequest, Message: "atleast one tag is required"},
		},
		{
			name: "should fail when got error while storing user tags",
			args: args{req: types.StoreTagRequest{Username: "Ankit", Publication: "AK", Tags: []string{"tag100"}}},
			mockDB: func() *handler.Application {
				tagStoreMock := mocks.NewUserTagStore(t)
				tagStoreMock.EXPECT().Store(mock.Anything, mock.Anything).Return(errors.New("db error"))

				m := model.Models{Tag: tagStoreMock}

				return handler.New(nil, &m)
			},
			wantRespBody: &response.Body{Status: http.StatusInternalServerError, Message: "error while storing user tag"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := tt.mockDB()

			handlerFunc := app.Store()
			got, gotErr := callEndpoint(t, &tt.args.req, handlerFunc, nil)

			if tt.wantRespBody != nil {
				assert.Nil(t, gotErr)
				assert.Equal(t, got.Status, tt.wantRespBody.Status)
				assert.Equal(t, got.Message, tt.wantRespBody.Message)
			}
		})
	}
}

func Test_Get(t *testing.T) {
	type args struct {
		urlParams map[string]string
	}

	tests := []struct {
		name         string
		args         args
		mockDB       func() *handler.Application
		wantResp     string
		wantRespBody *response.Body
	}{
		{
			name: "success",
			args: args{urlParams: map[string]string{"username": "Ankit", "publication": "AK"}},
			mockDB: func() *handler.Application {
				tagStoreMock := mocks.NewUserTagStore(t)
				tagStoreMock.EXPECT().Get(mock.Anything, mock.Anything).Return([]string{}, nil)

				m := model.Models{
					Tag: tagStoreMock,
				}

				return handler.New(nil, &m)
			},
			// wantResp:     handler.ArticleResponse{ID: 1},
			wantRespBody: &response.Body{Status: http.StatusOK, Message: ""},
		},
		{
			name: "should fail when invalid request is passed - empty username",
			args: args{urlParams: map[string]string{"username": "", "publication": "AK"}},
			mockDB: func() *handler.Application {
				m := model.Models{}

				return handler.New(nil, &m)
			},
			wantRespBody: &response.Body{Status: http.StatusBadRequest, Message: "username field is required"},
		},
		{
			name: "should fail when invalid request is passed - empty publication",
			args: args{urlParams: map[string]string{"username": "Ankit", "publication": ""}},
			mockDB: func() *handler.Application {
				m := model.Models{}

				return handler.New(nil, &m)
			},
			wantRespBody: &response.Body{Status: http.StatusBadRequest, Message: "publication field is required"},
		},
		{
			name: "should fail when got error while storing user tags",
			args: args{urlParams: map[string]string{"username": "Ankit", "publication": "AK"}},
			mockDB: func() *handler.Application {
				tagStoreMock := mocks.NewUserTagStore(t)
				tagStoreMock.EXPECT().Get(mock.Anything, mock.Anything).Return([]string{}, errors.New("db error"))

				m := model.Models{Tag: tagStoreMock}

				return handler.New(nil, &m)
			},
			wantRespBody: &response.Body{Status: http.StatusInternalServerError, Message: "error while fetching user tags"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := tt.mockDB()

			handlerFunc := app.Get()

			got, gotErr := callEndpoint(t, nil, handlerFunc, tt.args.urlParams)
			if tt.wantRespBody != nil {
				assert.Nil(t, gotErr)
				assert.Equal(t, got.Status, tt.wantRespBody.Status)
				assert.Equal(t, got.Message, tt.wantRespBody.Message)
			}
		})
	}
}

// callEndpoint creates a request and make a http call
func callEndpoint(t *testing.T, req *types.StoreTagRequest, handlerFunc http.HandlerFunc, urlParams map[string]string) (*response.Body, error) {
	w := httptest.NewRecorder()

	rawReq, _ := json.Marshal(req)

	// create a request
	r, err := http.NewRequest(mock.Anything, "/tag", bytes.NewBuffer(rawReq))
	if err != nil {
		t.Fatal(err)
	}

	// appends a urlParams at the end of route
	r = setURLParams(r, urlParams)

	// server http call
	handlerFunc.ServeHTTP(w, r)

	resp := response.Body{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Errorf("error unmarshalling response : %v", err)
	}

	return &resp, nil
}

// setURLParams appends a urlParams at the end of route
func setURLParams(req *http.Request, urlParams map[string]string) *http.Request {
	if len(urlParams) > 0 {
		routeContext := chi.NewRouteContext()
		for k, v := range urlParams {
			routeContext.URLParams.Add(k, v)
		}

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))
	}

	return req
}
