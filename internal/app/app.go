package app

import (
	"database/sql"
	"fmt"
	"github.com/nickemma/internal/api"
	"github.com/nickemma/internal/store"
	"github.com/nickemma/migration"
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
	// database connections
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}
	// migrations run and check
	err = store.MigrateFs(pgDB, migration.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
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
