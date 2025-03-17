package main

import (
	"github.com/nickemma/internal/app"
	"net/http"
	"time"
)

func main() {
	app, err := app.NewApplication()

	if err != nil {
		panic(err)
	}
	app.Logger.Println("We are running our application!")

	server := &http.Server{
		Addr:         "8080",
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}
