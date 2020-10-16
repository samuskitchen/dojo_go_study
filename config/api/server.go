package api

import (
	"dojo_go_study/config/database"
	"github.com/go-chi/chi/middleware"
	"log"

	"net/http"
	"time"

	"github.com/go-chi/chi"
)

// Server is a base server configuration.
type Server struct {
	handler http.Server
}

// ServeHTTP implements the http.Handler interface for the Server type.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.Handler.ServeHTTP(w, r)
}

// NewApplication initialized a new server with configuration.
func NewApplication(port string, conn *database.Data) *Server {

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Mount("/api/v1", New(conn))

	server := Server{handler: http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}}

	return &server
}

// Close server resources.
func (s *Server) Close() error {
	// TODO: add resource closure.
	return s.handler.Close()
}

// Start the server.
func (s *Server) Start() {
	log.Printf("handler running on http://localhost%s", s.handler.Addr)
	log.Fatal(s.handler.ListenAndServe())
}
