package main

import (
	"fmt"
	"sublinks/sublinks-federation/internal/db"
	"sublinks/sublinks-federation/internal/http"
	"sublinks/sublinks-federation/internal/log"
	"sublinks/sublinks-federation/internal/queue"

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
	defer conn.Close()
	conn.RunMigrations()

	q := queue.NewQueue(logger)
	err = q.Connect()
	if err != nil {
		logger.Fatal("failed connecting to queue service", err)
	}
	defer q.Close()
	err = q.CreateProducer("backend")
	if err != nil {
		logger.Fatal("failed creating producer", err)
	}
	err = q.CreateConsumer("federation")
	if err != nil {
		logger.Fatal("failed creating consumer", err)
	}
	q.StartConsumer("federation")
	config := http.ServerConfig{
		Logger:   logger,
		Database: conn,
		Queue:    q,
	}
	s := http.NewServer(config)
	s.RunServer()
}
