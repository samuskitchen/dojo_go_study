package test

import (
	"dojo_go_study/config/api"
	"dojo_go_study/config/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
	"time"
)

// NewServerTest initialized a newRoutes server with configuration.
func NewServerTest(port string, conn *database.Data) *api.Server {

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Mount("/api/v1", api.Routes(conn))

	s := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &api.Server{Server: s}
}
