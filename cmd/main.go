package main

import (
	"article-tag/internal/database"
	"article-tag/internal/handler"
	"article-tag/internal/model"
	"article-tag/internal/routes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

var app *handler.Application

// init
func init() {
	// initialize database
	db, err := database.InitDB()
	if err != nil {
		panic(err)
	}

	models := model.NewModel(db)

	// check and create table
	err = checkAndCreateTable(&models)
	if err != nil {
		panic(err)
	}

	app = handler.New(db, &models)
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
