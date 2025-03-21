package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/nickemma/internal/app"
)

func SetUpRoute(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)
	r.Get("/workouts/{id}", app.WorkoutHandler.HandlerGetWorkoutByID)

	r.Post("/workouts", app.WorkoutHandler.HandlerCreateWorkout)
	r.Put("/workouts/{id}", app.WorkoutHandler.HandleUpdateWorkoutById)
	r.Delete("/workouts/{id}", app.WorkoutHandler.HandleDeleteWorkoutById)

	return r
}
