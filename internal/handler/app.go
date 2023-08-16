package handler

import (
	"article-tag/internal/model"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Application struct {
	db    *dynamodb.Client
	model model.Models
}

// New
func New(db *dynamodb.Client, models *model.Models) *Application {
	return &Application{
		db:    db,
		model: *models,
	}
}
