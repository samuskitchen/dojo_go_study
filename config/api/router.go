package api

import (
	"dojo_go_study/config/database"
	"dojo_go_study/handler"
	"github.com/go-chi/chi"
	"net/http"
)

// Routes returns the API V1 Handler with configuration.
func Routes(conn *database.Data) http.Handler {
	router := chi.NewRouter()

	ur := handler.NewUserHandler(conn)
	router.Mount("/users", routesUser(ur))

	return router
}

// routesUser returns user router with each endpoint.
func routesUser(handler *handler.UserRouter) http.Handler {
	router := chi.NewRouter()

	router.Get("/", handler.GetAllUsersHandler)
	router.Get("/{id}", handler.GetOneHandler)
	router.Post("/", handler.CreateHandler)
	router.Put("/{id}", handler.UpdateHandler)
	router.Delete("/{id}", handler.DeleteHandler)

	return router
}