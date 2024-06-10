package http

import (
	"context"
	"net/http"
	"os"
	"sublinks/sublinks-federation/internal/db"
	"sublinks/sublinks-federation/internal/log"
	"sublinks/sublinks-federation/internal/queue"
	"sublinks/sublinks-federation/internal/service"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	*mux.Router
	log.Logger
	db.Database
	queue.Queue
	ServiceManager *service.ServiceManager
	srv            *http.Server
}

type ServerConfig struct {
	log.Logger
	db.Database
	queue.Queue
	ServiceManager *service.ServiceManager
}

func NewServer(config ServerConfig) *Server {
	r := mux.NewRouter()

	return &Server{
		Router:         r,
		Logger:         config.Logger,
		Database:       config.Database,
		Queue:          config.Queue,
		ServiceManager: config.ServiceManager,
		srv: &http.Server{
			Addr: os.Getenv("LISTEN_ADDR"),
			// Good practice to set timeouts to avoid Slowloris attacks.
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			// pass embed of Server for *mux
			Handler: r,
		},
	}
}

func (server *Server) RunServer(ctx context.Context) {

	server.SetupInternalRoutes()
	server.SetupUserRoutes()
	server.SetupPostRoutes()
	server.SetupApubRoutes()
	server.SetupActivityRoutes()
	server.Router.NotFoundHandler = http.HandlerFunc(server.notFound)
	server.Router.MethodNotAllowedHandler = http.HandlerFunc(server.notAllowedMethod)
	server.Router.Use(server.logMiddleware)

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		server.Logger.Info("Starting server")
		if err := server.srv.ListenAndServe(); err != nil {
			server.Logger.Error("Error starting server", err)
		}
	}()

}

func (server *Server) Shutdown(ctx context.Context) error {
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	err := server.srv.Shutdown(ctx)
	if err != nil {
		server.Logger.Error("Error shutting down server", err)
		return err
	}
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	server.Logger.Info("shutting down")
	return nil
}
