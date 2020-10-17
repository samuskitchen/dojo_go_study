package api

import (
	"dojo_go_study/handler"
	"github.com/go-chi/chi"
	"net/http"
)


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