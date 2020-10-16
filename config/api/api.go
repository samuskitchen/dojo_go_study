package api

import (
	"dojo_go_study/config/database"
	"dojo_go_study/handler"
	"net/http"

	"github.com/go-chi/chi"
)

// New returns the API V1 Handler with configuration.
func New(conn *database.Data) http.Handler {
	router := chi.NewRouter()

	ur := handler.NewUserHandler(conn)
	router.Mount("/users", RoutesUser(ur))

	return router
}
