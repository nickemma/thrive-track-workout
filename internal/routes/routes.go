package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/nickemma/internal/app"
)

func SetUpRoute(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)

	return r
}
