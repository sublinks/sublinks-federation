package main

import (
	"fmt"
	"sublinks/sublinks-federation/internal/db"
	"sublinks/sublinks-federation/internal/http"
	"sublinks/sublinks-federation/internal/log"
	"sublinks/sublinks-federation/internal/queue"
	"sublinks/sublinks-federation/internal/service"
	"sublinks/sublinks-federation/internal/service/actors"
	"sublinks/sublinks-federation/internal/service/comments"
	"sublinks/sublinks-federation/internal/service/posts"

	"github.com/joho/godotenv"
)

func main() {
	// bootstrap logger
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
		actors.NewUserService(conn),
		actors.NewCommunityService(conn),
		posts.NewPostService(conn),
		comments.NewCommentService(conn),
	)
	q.Run(serviceManager)
	config := http.ServerConfig{
		Logger:         logger,
		Queue:          q,
		ServiceManager: serviceManager,
	}
	s := http.NewServer(config)
	s.RunServer()
}
