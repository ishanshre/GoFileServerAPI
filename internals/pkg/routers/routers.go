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
	r.Post("/api/v1/register", h.UserRegister)
	r.Post("/api/v1/login", h.UserLogin)
	r.Group(func(admin chi.Router) {
		admin.Use(m.JwtAuth)
		admin.Use(m.CheckAdmin)
		admin.Get("/api/v1/admin/users", h.GetUsers)
		admin.Get("/api/v1/admin/users/{username}", h.GetUser)
		admin.Delete("/api/v1/admin/users/{username}", h.AdminDeleteUser)
	})
	r.Group(func(user chi.Router) {
		user.Use(m.JwtAuth)
		user.Post("/api/v1/logout", h.UserLogout)
		user.Get("/api/v1/me", h.GetMe)
		user.Delete("/api/v1/me", h.DeleteMe)
	})
	return r
}
