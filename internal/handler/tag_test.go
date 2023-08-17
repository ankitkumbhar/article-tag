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
			args: args{req: types.StoreTagRequest{Username: "Test", Publication: "TS", Tags: []string{"tag100"}}},
			mockDB: func() *handler.Application {
				tagStoreMock := mocks.NewUserTagStore(t)
				tagStoreMock.EXPECT().DescribeTable(mock.Anything).Return(nil)
				tagStoreMock.EXPECT().Store(mock.Anything, mock.Anything).Return(nil)

				m := model.Models{
					Tag: tagStoreMock,
				}

				return handler.New(nil, &m)
			},
			wantRespBody: &response.Body{Status: http.StatusCreated, Message: ""},
		},
		{
			name: "should fail when invalid request is passed - empty username",
			args: args{req: types.StoreTagRequest{Username: "", Publication: "TS", Tags: []string{"tag100"}}},
			mockDB: func() *handler.Application {
				m := model.Models{}

				return handler.New(nil, &m)
			},
			wantRespBody: &response.Body{Status: http.StatusBadRequest, Message: "username field is required"},
		},
		{
			name: "should fail when invalid request is passed - empty publication",
			args: args{req: types.StoreTagRequest{Username: "Test", Publication: "", Tags: []string{"tag100"}}},
			mockDB: func() *handler.Application {
				m := model.Models{}

				return handler.New(nil, &m)
			},
			wantRespBody: &response.Body{Status: http.StatusBadRequest, Message: "publication field is required"},
		},
		{
			name: "should fail when invalid request is passed - empty tags",
			args: args{req: types.StoreTagRequest{Username: "Test", Publication: "TS", Tags: []string{}}},
			mockDB: func() *handler.Application {
				m := model.Models{}

				return handler.New(nil, &m)
			},
			wantRespBody: &response.Body{Status: http.StatusBadRequest, Message: "atleast one tag is required"},
		},
		{
			name: "should fail when got error while storing user tags",
			args: args{req: types.StoreTagRequest{Username: "Test", Publication: "TS", Tags: []string{"tag101"}}},
			mockDB: func() *handler.Application {
				tagStoreMock := mocks.NewUserTagStore(t)
				tagStoreMock.EXPECT().DescribeTable(mock.Anything).Return(nil)
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

			rawReq, _ := json.Marshal(tt.args.req)
			got, gotErr := callEndpoint(t, rawReq, handlerFunc, nil)

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
			args: args{urlParams: map[string]string{"username": "Test", "publication": "TS"}},
			mockDB: func() *handler.Application {
				tagStoreMock := mocks.NewUserTagStore(t)
				tagStoreMock.EXPECT().Get(mock.Anything, mock.Anything).Return([]string{}, nil)

				m := model.Models{
					Tag: tagStoreMock,
				}

				return handler.New(nil, &m)
			},
			wantRespBody: &response.Body{Status: http.StatusOK, Message: ""},
		},
		{
			name: "should fail when invalid request is passed - empty username",
			args: args{urlParams: map[string]string{"username": "", "publication": "TS"}},
			mockDB: func() *handler.Application {
				m := model.Models{}

				return handler.New(nil, &m)
			},
			wantRespBody: &response.Body{Status: http.StatusBadRequest, Message: "username field is required"},
		},
		{
			name: "should fail when invalid request is passed - empty publication",
			args: args{urlParams: map[string]string{"username": "Test", "publication": ""}},
			mockDB: func() *handler.Application {
				m := model.Models{}

				return handler.New(nil, &m)
			},
			wantRespBody: &response.Body{Status: http.StatusBadRequest, Message: "publication field is required"},
		},
		{
			name: "should fail when got error while storing user tags",
			args: args{urlParams: map[string]string{"username": "Test", "publication": "TS"}},
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

func Test_Delete(t *testing.T) {
	type args struct {
		req types.DeleteTagRequest
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
			args: args{types.DeleteTagRequest{Username: "Test", Publication: "TS", Tags: []string{"tag101"}}},
			mockDB: func() *handler.Application {
				tagStoreMock := mocks.NewUserTagStore(t)
				tagStoreMock.EXPECT().Delete(mock.Anything, mock.Anything).Return(nil)

				m := model.Models{
					Tag: tagStoreMock,
				}

				return handler.New(nil, &m)
			},
			wantRespBody: &response.Body{Status: http.StatusOK, Message: ""},
		},
		{
			name: "should fail when invalid request is passed - empty username",
			args: args{types.DeleteTagRequest{Username: "", Publication: "TS", Tags: []string{"tag101"}}},
			mockDB: func() *handler.Application {
				m := model.Models{}

				return handler.New(nil, &m)
			},
			wantRespBody: &response.Body{Status: http.StatusBadRequest, Message: "username field is required"},
		},
		{
			name: "should fail when invalid request is passed - empty publication",
			args: args{types.DeleteTagRequest{Username: "Test", Publication: "", Tags: []string{"tag101"}}},
			mockDB: func() *handler.Application {
				m := model.Models{}

				return handler.New(nil, &m)
			},
			wantRespBody: &response.Body{Status: http.StatusBadRequest, Message: "publication field is required"},
		},
		{
			name: "Should fail when receive error from database while deleting userTag",
			args: args{types.DeleteTagRequest{Username: "Test", Publication: "TS", Tags: []string{"tag101"}}},
			mockDB: func() *handler.Application {
				tagStoreMock := mocks.NewUserTagStore(t)
				tagStoreMock.EXPECT().Delete(mock.Anything, mock.Anything).Return(errors.New("db error"))

				m := model.Models{
					Tag: tagStoreMock,
				}

				return handler.New(nil, &m)
			},
			wantRespBody: &response.Body{Status: http.StatusInternalServerError, Message: "error while deleting user followed tags"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := tt.mockDB()

			handlerFunc := app.Delete()

			rawReq, _ := json.Marshal(tt.args.req)

			got, gotErr := callEndpoint(t, rawReq, handlerFunc, nil)
			if tt.wantRespBody != nil {
				assert.Nil(t, gotErr)
				assert.Equal(t, got.Status, tt.wantRespBody.Status)
				assert.Equal(t, got.Message, tt.wantRespBody.Message)
			}
		})
	}
}

func Test_PopularTags(t *testing.T) {
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
			args: args{urlParams: map[string]string{"username": "Test", "publication": "TS"}},
			mockDB: func() *handler.Application {
				tagStoreMock := mocks.NewUserTagStore(t)
				tagStoreMock.EXPECT().GetPopularTags(mock.Anything, mock.Anything).Return([]string{"tag101"}, nil)

				m := model.Models{
					Tag: tagStoreMock,
				}

				return handler.New(nil, &m)
			},
			wantRespBody: &response.Body{Status: http.StatusOK, Message: ""},
		},
		{
			name: "should fail when invalid request is passed - empty username",
			args: args{urlParams: map[string]string{"username": "", "publication": "TS"}},
			mockDB: func() *handler.Application {
				m := model.Models{}

				return handler.New(nil, &m)
			},
			wantRespBody: &response.Body{Status: http.StatusBadRequest, Message: "username field is required"},
		},
		{
			name: "should fail when invalid request is passed - empty publication",
			args: args{urlParams: map[string]string{"username": "Test", "publication": ""}},
			mockDB: func() *handler.Application {
				m := model.Models{}

				return handler.New(nil, &m)
			},
			wantRespBody: &response.Body{Status: http.StatusBadRequest, Message: "publication field is required"},
		},
		{
			name: "Should fail when receive error from database while fetching popular userTag",
			args: args{urlParams: map[string]string{"username": "Test", "publication": "TS"}},
			mockDB: func() *handler.Application {
				tagStoreMock := mocks.NewUserTagStore(t)
				tagStoreMock.EXPECT().GetPopularTags(mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

				m := model.Models{
					Tag: tagStoreMock,
				}

				return handler.New(nil, &m)
			},
			wantRespBody: &response.Body{Status: http.StatusInternalServerError, Message: "error while fetching popular tags"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := tt.mockDB()

			handlerFunc := app.PopularTag()

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
func callEndpoint(t *testing.T, rawReq []byte, handlerFunc http.HandlerFunc, urlParams map[string]string) (*response.Body, error) {
	w := httptest.NewRecorder()

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
