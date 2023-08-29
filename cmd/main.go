package main

import (
	"article-tag/internal/database"
	"article-tag/internal/handler"
	"article-tag/internal/model"
	"article-tag/internal/routes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var app *handler.Application

// init
func init() {
	// initialize database
	db, err := database.InitDB()
	if err != nil {
		panic(err)
	}

	// initialize logger
	logger := initLogger()

	models := model.NewModel(db, logger)

	// check and create table
	err = checkAndCreateTable(&models)
	if err != nil {
		panic(err)
	}

	app = handler.New(db, &models, logger)
}

func initLogger() *zap.Logger {
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

	// add timestamp in the log
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger := zap.Must(cfg.Build())

	return logger
}

// checkAndCreateTable
func checkAndCreateTable(models *model.Models) error {
	// wait for some time until docker is up
	time.Sleep(time.Second * 2)

	// check table exists
	err := models.Tag.DescribeTable(context.TODO())
	if err != nil {
		err = models.Tag.CreateTable(context.TODO())
		if err != nil {
			log.Println("error creating table : ", err)

			return err
		}
	}

	return nil
}

func main() {
	r := routes.InitRouter(app)

	port := "8080"

	log.Default().Println("starting server on port :", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), r); err != nil {
		log.Fatalf("error starting server on port : %v", port)
	}
}
