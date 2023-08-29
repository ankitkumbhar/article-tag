package model

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.uber.org/zap"
)

type UserTagStore interface {
	DescribeTable(ctx context.Context) error
	CreateTable(ctx context.Context) error
	Store(ctx context.Context, username, publication, tagID, tagName string) error
	Get(ctx context.Context, username, publication, order string) ([]*UserTag, error)
	Delete(ctx context.Context, username, publication, tagID, tagName string) error
	GetPopularTags(ctx context.Context, username, publication string) ([]string, error)
}

type UserTag struct {
	PK          string
	SK          string
	TagID       string
	TagName     string
	CreatedAt   string
	Username    string
	Publication string
}

// ExclusiveStartKey
type ExclusiveStartKey struct {
	SK       string
	TagCount string
}

type Models struct {
	Tag UserTagStore
}

func NewModel(db *dynamodb.Client, logger *zap.Logger) Models {
	return Models{
		Tag: NewTag(db, logger),
	}
}
