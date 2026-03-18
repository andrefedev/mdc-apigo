package okhttpx

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type AppRouterDeps struct {
	AppHandler  http.Handler
	AuthHandler http.Handler
	UserHandler http.Handler

	ReadyHandler http.Handler
}

func NewAppRouter(deps AppRouterDeps) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.CleanPath)

	r.Get("/healthz", Healthz)
	r.Get("/readyz", deps.ReadyHandler.ServeHTTP)

	r.Route("/v1", func(v1 chi.Router) {
		v1.Mount("/app", deps.AppHandler)
		v1.Mount("/auth", deps.AuthHandler)
		v1.Mount("/users", deps.UserHandler)
	})

	return r
}
