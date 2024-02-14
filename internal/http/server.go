package http

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"sublinks/sublinks-federation/internal/db"
	"sublinks/sublinks-federation/internal/log"
	"sublinks/sublinks-federation/internal/queue"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	*mux.Router
	log.Logger
	*db.Database
	queue.Queue
}

type ServerConfig struct {
	log.Logger
	*db.Database
	queue.Queue
}

func NewServer(config ServerConfig) *Server {
	r := mux.NewRouter()

	return &Server{
		Router:   r,
		Logger:   config.Logger,
		Database: config.Database,
		Queue:    config.Queue,
	}
}

func (server *Server) RunServer() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	server.SetupUserRoutes()
	server.SetupPostRoutes()
	server.SetupApubRoutes()
	server.SetupActivityRoutes()
	server.Router.NotFoundHandler = http.HandlerFunc(server.notFound)
	server.Router.MethodNotAllowedHandler = http.HandlerFunc(server.notAllowedMethod)
	server.Router.Use(server.logMiddleware)

	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		// pass embed of Server for *mux
		Handler: server.Router,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		server.Logger.Info("Starting server")
		if err := srv.ListenAndServe(); err != nil {
			server.Logger.Error("Error starting server", err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	err := srv.Shutdown(ctx)
	if err != nil {
		server.Logger.Error("Error shutting down server", err)
	}
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	server.Logger.Info("shutting down")
}
