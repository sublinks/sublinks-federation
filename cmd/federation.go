package main

import (
	"fmt"
	"os"
	"sublinks/sublinks-federation/internal/db"
	"sublinks/sublinks-federation/internal/http"
	"sublinks/sublinks-federation/internal/log"
	"sublinks/sublinks-federation/internal/queue"

	"github.com/joho/godotenv"
)

func main() {
	// Load connection string from .env file
	err := godotenv.Load()
	if err != nil {
		log.Warn(fmt.Sprintf("failed to load env, %v", err))
	}
	conn, err := db.Connect()
	if err != nil {
		log.Fatal("failed connecting to db", err)
	}
	defer conn.Close()
	db.RunMigrations(conn)

	mqConnection, err := queue.Connect()
	if err != nil {
		log.Fatal("failed connecting to queue service", err)
	}
	defer mqConnection.Close()
	producer, err := queue.CreateProducer(mqConnection, "backend")
	if err != nil {
		log.Fatal("failed creating producer", err)
	}
	defer producer.Close()
	messages, err := queue.CreateConsumer(mqConnection, "federation")
	if err != nil {
		log.Fatal("failed creating consumer", err)
	}
	go func() {
		for message := range messages {
			log.Debug(fmt.Sprintf(" > Received message: %s\n", message.Body))
		}
	}()

	http.RunServer()
	os.Exit(0)
}
