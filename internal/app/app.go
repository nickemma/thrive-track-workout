package app

import (
	"database/sql"
	"fmt"
	"github.com/nickemma/internal/api"
	"github.com/nickemma/internal/store"
	"log"
	"net/http"
	"os"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	DB             *sql.DB
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

	// database connections
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}
	// Store goes here

	// Handlers goes here
	workoutHandler := api.NewWorkoutHandler()

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		DB:             pgDB,
	}

	return app, nil
}

func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available and ok\n")

}
