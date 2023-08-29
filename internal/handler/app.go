package handler

import (
	"article-tag/internal/model"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Application struct {
	db       *dynamodb.Client
	model    model.Models
	validate *validator.Validate
	logger   *zap.Logger
}

// New
func New(db *dynamodb.Client, models *model.Models, logger *zap.Logger) *Application {
	// return app object
	return &Application{
		db:       db,
		model:    *models,
		validate: validator.New(),
		logger:   logger,
	}
}

// GetLogger
func GetLogger(app *Application) *zap.Logger {
	// return logger object
	return app.logger
}
