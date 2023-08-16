package main

import (
	"article-tag/internal/database"
	"article-tag/internal/handler"
	"article-tag/internal/model"
	"article-tag/internal/routes"
	"fmt"
	"log"
	"net/http"
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

	app = handler.New(db, &models)
}

func main() {
	r := routes.InitRouter(app)

	port := "8080"

	log.Default().Println("starting server on port :", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), r); err != nil {
		log.Fatalf("error starting server on port : %v", port)
	}
}
