package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/ishanshre/GoFileServerAPI/internals/pkg/handlers"
	"github.com/ishanshre/GoFileServerAPI/internals/pkg/middlewares"
)

func NewRouter(h handlers.Handlers, m middlewares.Middlewares) http.Handler {
	r := chi.NewRouter()
	r.Use(cors.Handler((cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})))
	r.Use(middleware.Logger)
	return r
}