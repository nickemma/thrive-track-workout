package main

import (
	"github.com/nickemma/internal/app"
	"github.com/nickemma/internal/routes"
	"net/http"
	"time"
)

func main() {
	app, err := app.NewApplication()

	if err != nil {
		panic(err)
	}
	app.Logger.Println("We are running our application!")

	r := routes.SetUpRoute(app)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}
