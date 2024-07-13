package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sublinks/sublinks-federation/internal/db"
	"sublinks/sublinks-federation/internal/http"
	"sublinks/sublinks-federation/internal/log"
	"sublinks/sublinks-federation/internal/queue"
	"sublinks/sublinks-federation/internal/service"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := &sync.WaitGroup{}

	// bootstrap logger
	logLevel, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		logLevel = "info"
	}
	log.SetGlobalLevel(logLevel)
	logger := log.NewLogger("main")

	// Load connection string from .env file
	err := godotenv.Load()
	if err != nil {
		logger.Warn(fmt.Sprintf("failed to load env, %v", err))
	}

	conn := db.NewDatabase()
	err = conn.Connect()
	if err != nil {
		logger.Fatal("failed connecting to db", err)
	}
	conn.RunMigrations()

	q := queue.NewQueue(logger)
	err = q.Connect()
	if err != nil {
		logger.Fatal("failed connecting to queue service", err)
	}
	defer q.Close()

	serviceManager := service.NewServiceManager(
		service.NewUserService(conn),
		service.NewCommunityService(conn),
		service.NewPostService(conn),
		service.NewCommentService(conn),
	)

	q.Run(ctx, serviceManager, wg)

	config := http.ServerConfig{
		Logger:         logger,
		Queue:          q,
		ServiceManager: serviceManager,
	}
	s := http.NewServer(config)
	go func() {
		s.RunServer(ctx)
	}()

	signalTermChan := make(chan os.Signal, 1)
	signal.Notify(signalTermChan, os.Interrupt)
	<-signalTermChan

	logger.Debug("shutting down gracefully")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), wait)
	defer shutdownCancel()

	cancel()
	wg.Wait()

	err = s.Shutdown(shutdownCtx)
	if err != nil {
		logger.Fatal("failed shutting down server", err)
	}
	logger.Info("shutdown complete")
}
