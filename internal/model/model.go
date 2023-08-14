package model

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type UserTagStore interface {
	Store(ctx context.Context, item *UserTag) error
	Get(ctx context.Context) ([]*UserTag, error)
}

type UserTag struct {
	SK   string
	PU   string
	Tag  string
	Tags []string
}

type Models struct {
	Tag UserTagStore
}

func NewModel(db *dynamodb.Client) Models {
	return Models{
		Tag: &tag{db: db},
	}
}
