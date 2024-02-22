package main

import (
	"fmt"
	"sublinks/sublinks-federation/internal/db"
	"sublinks/sublinks-federation/internal/http"
	"sublinks/sublinks-federation/internal/log"
	"sublinks/sublinks-federation/internal/queue"
	"sublinks/sublinks-federation/internal/repository"
	"sublinks/sublinks-federation/internal/worker"

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
	processActors(logger, conn, q)
	processPosts(logger, conn, q)
	config := http.ServerConfig{
		Logger:   logger,
		Database: conn,
		Queue:    q,
	}
	s := http.NewServer(config)
	s.RunServer()
}

func processActors(logger *log.Log, conn db.Database, q queue.Queue) {
	actorCQ := queue.ConsumerQueue{
		QueueName:  "actor_create_queue",
		Exchange:   "federation",
		RoutingKey: "actor.create",
	}

	aw := worker.ActorWorker{
		Logger:     logger,
		Repository: repository.NewRepository(conn),
	}

	err := q.StartConsumer(actorCQ, &aw)
	if err != nil {
		logger.Fatal("failed starting actor consumer", err)
	}
}

func processPosts(logger *log.Log, conn db.Database, q queue.Queue) {
	postCQ := queue.ConsumerQueue{
		QueueName:  "post_queue",
		Exchange:   "federation",
		RoutingKey: "post.create",
	}

	aw := worker.PostWorker{
		Logger:     logger,
		Repository: repository.NewRepository(conn),
	}

	err := q.StartConsumer(postCQ, &aw)
	if err != nil {
		logger.Fatal("failed starting post consumer", err)
	}
}
