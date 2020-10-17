package api

import (
	"context"
	"dojo_go_study/config/database"
	"dojo_go_study/handler"
	"github.com/go-chi/chi/middleware"
	"log"
	"os"
	"os/signal"

	"net/http"
	"time"

	"github.com/go-chi/chi"
)

// server is a base server configuration.
type server struct {
	*http.Server
}

// ServeHTTP implements the http.Handler interface for the server type.
func (srv *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.Handler.ServeHTTP(w, r)
}

// newServer initialized a newRoutes server with configuration.
func newServer(port string, conn *database.Data) *server {

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Mount("/api/v1", newRoutes(conn))

	s := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &server{s}
}

// newRoutes returns the API V1 Handler with configuration.
func newRoutes(conn *database.Data) http.Handler {
	router := chi.NewRouter()

	ur := handler.NewUserHandler(conn)
	router.Mount("/users", routesUser(ur))

	return router
}

// Start the server.
func (srv *server) Start() {
	log.Println("starting API cmd")

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not listen on %s rv due to %s rv", srv.Addr, err.Error())
		}
	}()
	log.Printf("cmd is ready to handle requests %s", srv.Addr)
	srv.gracefulShutdown()
}

func (srv *server) gracefulShutdown() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	log.Printf("cmd is shutting down %s", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	srv.SetKeepAlivesEnabled(false)
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("could not gracefully shutdown the cmd %s", err.Error())
	}
	log.Printf("cmd stopped")
}