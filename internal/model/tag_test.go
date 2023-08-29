package model_test

import (
	"article-tag/internal/mocks"
	"article-tag/internal/model"
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func testSuite() *zap.Logger {
	rawJSON := []byte(`{
		"level": "debug",
		"encoding": "json",
		"outputPaths": ["stdout", "/tmp/logs"],
		"errorOutputPaths": ["stderr"],
		"encoderConfig": {
		  "messageKey": "message",
		  "levelKey": "level",
		  "levelEncoder": "lowercase"
		}
	  }`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}

	return zap.Must(cfg.Build())
}

func Test_Describe(t *testing.T) {
	log := testSuite()

	type args struct {
		item model.UserTag
	}

	tests := []struct {
		name    string
		args    args
		mockDB  func() model.Models
		wantErr error
	}{
		{
			name: "success",
			args: args{item: model.UserTag{Username: "Mock username"}},
			mockDB: func() model.Models {
				models := model.NewModel(nil, log)

				dmock := mocks.NewDynamoAPI(t)
				dmock.EXPECT().DescribeTable(mock.Anything, mock.Anything).Return(&dynamodb.DescribeTableOutput{}, nil)
				models.Tag = model.NewTag(dmock, log)

				return models
			},
			wantErr: nil,
		},
		{
			name: "Should fail when received error in describe table",
			args: args{item: model.UserTag{Username: "Mock username"}},
			mockDB: func() model.Models {
				models := model.NewModel(nil, log)

				dmock := mocks.NewDynamoAPI(t)
				dmock.EXPECT().DescribeTable(mock.Anything, mock.Anything).Return(&dynamodb.DescribeTableOutput{}, errors.New("mock error"))
				models.Tag = model.NewTag(dmock, log)

				return models
			},
			wantErr: errors.New("mock error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.mockDB()

			// call model function
			err := a.Tag.DescribeTable(context.TODO())

			if tt.wantErr == nil {
				assert.Equal(t, tt.wantErr, err)
			}

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			}
		})
	}
}

func Test_CreateTable(t *testing.T) {
	log := testSuite()

	type args struct {
		item model.UserTag
	}

	tests := []struct {
		name    string
		args    args
		mockDB  func() model.Models
		wantErr error
	}{
		{
			name: "success",
			args: args{item: model.UserTag{Username: "Mock username"}},
			mockDB: func() model.Models {
				models := model.NewModel(nil, log)

				dmock := mocks.NewDynamoAPI(t)
				dmock.EXPECT().CreateTable(mock.Anything, mock.Anything).Return(&dynamodb.CreateTableOutput{}, nil)
				models.Tag = model.NewTag(dmock, log)

				return models
			},
			wantErr: nil,
		},
		{
			name: "Should fail when received error in describe table",
			args: args{item: model.UserTag{Username: "Mock username"}},
			mockDB: func() model.Models {
				models := model.NewModel(nil, log)

				dmock := mocks.NewDynamoAPI(t)
				dmock.EXPECT().CreateTable(mock.Anything, mock.Anything).Return(&dynamodb.CreateTableOutput{}, errors.New("mock error"))
				models.Tag = model.NewTag(dmock, log)

				return models
			},
			wantErr: errors.New("mock error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.mockDB()

			// call model function
			err := a.Tag.CreateTable(context.TODO())

			if tt.wantErr == nil {
				assert.Equal(t, tt.wantErr, err)
			}

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			}
		})
	}
}

func Test_Store(t *testing.T) {
	log := testSuite()

	type args struct {
		item model.UserTag
	}

	tests := []struct {
		name    string
		args    args
		mockDB  func() model.Models
		wantErr error
	}{
		{
			name: "success",
			args: args{item: model.UserTag{Username: "Mock username"}},
			mockDB: func() model.Models {
				models := model.NewModel(nil, log)

				dmock := mocks.NewDynamoAPI(t)
				dmock.EXPECT().PutItem(mock.Anything, mock.Anything).Return(&dynamodb.PutItemOutput{}, nil)
				dmock.EXPECT().UpdateItem(mock.Anything, mock.Anything).Return(&dynamodb.UpdateItemOutput{}, nil)
				models.Tag = model.NewTag(dmock, log)

				return models
			},
			wantErr: nil,
		},
		{
			name: "Should fail when received error in putItem call",
			args: args{item: model.UserTag{Username: "Mock username"}},
			mockDB: func() model.Models {
				models := model.NewModel(nil, log)

				dmock := mocks.NewDynamoAPI(t)
				dmock.EXPECT().PutItem(mock.Anything, mock.Anything).Return(nil, errors.New("mock error"))
				models.Tag = model.NewTag(dmock, log)

				return models
			},
			wantErr: errors.New("mock error"),
		},
		{
			name: "Should fail when received error in updateItem call",
			args: args{item: model.UserTag{Username: "Mock username"}},
			mockDB: func() model.Models {
				models := model.NewModel(nil, log)

				dmock := mocks.NewDynamoAPI(t)
				dmock.EXPECT().PutItem(mock.Anything, mock.Anything).Return(&dynamodb.PutItemOutput{}, nil)
				dmock.EXPECT().UpdateItem(mock.Anything, mock.Anything).Return(&dynamodb.UpdateItemOutput{}, errors.New("mock error"))
				models.Tag = model.NewTag(dmock, log)

				return models
			},
			wantErr: errors.New("mock error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.mockDB()

			// call model function
			err := a.Tag.Store(context.TODO(), tt.args.item.Username, tt.args.item.Publication, tt.args.item.TagID, tt.args.item.TagName)

			if tt.wantErr == nil {
				assert.Equal(t, tt.wantErr, err)
			}

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			}
		})
	}
}

func Test_Get(t *testing.T) {
	log := testSuite()

	type args struct {
		item  model.UserTag
		order string
	}

	tests := []struct {
		name    string
		args    args
		mockDB  func() model.Models
		wantErr error
		want    []*model.UserTag
	}{
		{
			name: "success",
			args: args{item: model.UserTag{Username: "Mock username"}},
			mockDB: func() model.Models {
				models := model.NewModel(nil, log)

				dmock := mocks.NewDynamoAPI(t)
				dmock.EXPECT().Query(mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{
					map[string]types.AttributeValue{
						"TagID":   &types.AttributeValueMemberS{Value: "1"},
						"TagName": &types.AttributeValueMemberS{Value: "tag1"},
					},
				}}, nil)
				models.Tag = model.NewTag(dmock, log)

				return models
			},
			wantErr: nil,
			want: []*model.UserTag{{
				TagID:   "1",
				TagName: "tag1",
			}},
		},
		{
			name: "Should fail when received error in query call",
			args: args{item: model.UserTag{Username: "Mock username"}},
			mockDB: func() model.Models {
				models := model.NewModel(nil, log)

				dmock := mocks.NewDynamoAPI(t)
				dmock.EXPECT().Query(mock.Anything, mock.Anything).Return(nil, errors.New("mock error"))
				models.Tag = model.NewTag(dmock, log)

				return models
			},
			wantErr: errors.New("mock error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.mockDB()

			// call model function
			tags, err := a.Tag.Get(context.TODO(), tt.args.item.Username, tt.args.item.Publication, tt.args.order)

			if tt.wantErr == nil {
				assert.Equal(t, tt.wantErr, err)
				assert.Equal(t, tt.want, tags)
			}

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			}
		})
	}
}

func Test_Delete(t *testing.T) {
	log := testSuite()

	type args struct {
		item model.UserTag
	}

	tests := []struct {
		name    string
		args    args
		mockDB  func() model.Models
		wantErr error
		want    []string
	}{
		{
			name: "success",
			args: args{item: model.UserTag{Username: "Mock username"}},
			mockDB: func() model.Models {
				models := model.NewModel(nil, log)

				dmock := mocks.NewDynamoAPI(t)
				dmock.EXPECT().DeleteItem(mock.Anything, mock.Anything).Return(&dynamodb.DeleteItemOutput{
					Attributes: map[string]types.AttributeValue{
						"PK": &types.AttributeValueMemberS{Value: "Test#AA"},
					},
				}, nil)

				dmock.EXPECT().UpdateItem(mock.Anything, mock.Anything).Return(&dynamodb.UpdateItemOutput{}, nil)
				models.Tag = model.NewTag(dmock, log)

				return models
			},
			wantErr: nil,
			want:    []string{"tag1"},
		},
		{
			name: "Should fail when received error in delete call",
			args: args{item: model.UserTag{Username: "Mock username"}},
			mockDB: func() model.Models {
				models := model.NewModel(nil, log)

				dmock := mocks.NewDynamoAPI(t)
				dmock.EXPECT().DeleteItem(mock.Anything, mock.Anything).Return(nil, errors.New("mock error"))
				models.Tag = model.NewTag(dmock, log)

				return models
			},
			wantErr: errors.New("mock error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.mockDB()

			// call model function
			err := a.Tag.Delete(context.TODO(), tt.args.item.Username, tt.args.item.Publication, tt.args.item.TagID, tt.args.item.TagName)

			if tt.wantErr == nil {
				assert.Equal(t, tt.wantErr, err)
			}

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			}
		})
	}
}

func Test_GetPopularTags(t *testing.T) {
	log := testSuite()

	type args struct {
		item model.UserTag
	}

	tests := []struct {
		name    string
		args    args
		mockDB  func() model.Models
		wantErr error
		want    []string
	}{
		{
			name: "success",
			args: args{item: model.UserTag{Username: "Mock username"}},
			mockDB: func() model.Models {
				models := model.NewModel(nil, log)

				dmock := mocks.NewDynamoAPI(t)
				dmock.EXPECT().Query(mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{
					map[string]types.AttributeValue{
						"TagName": &types.AttributeValueMemberS{Value: "tag1"},
					},
				}}, nil).Once()

				dmock.EXPECT().Query(mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{
					map[string]types.AttributeValue{
						"TagName": &types.AttributeValueMemberS{Value: "tag101"},
						"SK":      &types.AttributeValueMemberS{Value: "1"},
					},
				}}, nil)
				models.Tag = model.NewTag(dmock, log)

				return models
			},
			// wantErr: nil,
			want: []string{"tag101"},
		},
		{
			name: "Should fail when received error in query call",
			args: args{item: model.UserTag{Username: "Mock username"}},
			mockDB: func() model.Models {
				models := model.NewModel(nil, log)

				dmock := mocks.NewDynamoAPI(t)
				dmock.EXPECT().Query(mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{
					map[string]types.AttributeValue{
						"TagName": &types.AttributeValueMemberS{Value: "tag1"},
					},
				}}, nil).Once()
				dmock.EXPECT().Query(mock.Anything, mock.Anything).Return(nil, errors.New("mock error")).Once()
				models.Tag = model.NewTag(dmock, log)

				return models
			},
			wantErr: errors.New("mock error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.mockDB()

			// call model function
			userTags, err := a.Tag.GetPopularTags(context.TODO(), tt.args.item.Username, tt.args.item.Publication)

			if tt.wantErr == nil {
				assert.Equal(t, tt.wantErr, err)
				assert.Equal(t, tt.want, userTags)
			}

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			}
		})
	}
}
