package model

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type tag struct {
	db *dynamodb.Client
}

func (t *tag) Store(ctx context.Context, item *UserTag) error {
	return nil
}

func (m *tag) Get(ctx context.Context) ([]*UserTag, error) {
	mu := []*UserTag{}

	return mu, nil
}
