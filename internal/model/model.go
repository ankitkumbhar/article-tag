package model

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type UserTagStore interface {
	DescribeTable(ctx context.Context) error
	CreateTable(ctx context.Context) error
	Store(ctx context.Context, item *UserTag) error
	Get(ctx context.Context, item *UserTag) ([]string, error)
	Delete(ctx context.Context, item *UserTag) error
	GetPopularTags(ctx context.Context, item *UserTag) ([]string, error)
}

type UserTag struct {
	Username    string
	SK          string
	Publication string
	Tag         string
	Tags        []string
}

type Models struct {
	Tag UserTagStore
}

func NewModel(db *dynamodb.Client) Models {
	return Models{
		Tag: NewTag(db),
	}
}
