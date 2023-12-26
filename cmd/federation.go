package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"sublinks/sublinks-federation/internal/db"
	"sublinks/sublinks-federation/internal/http"
	"sublinks/sublinks-federation/internal/log"
	"sublinks/sublinks-federation/internal/queue"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
)

func SetupCloseHandler(queueService *amqp091.Channel, producer *amqp091.Connection, db *sql.DB) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Info("Ctrl+C pressed in Terminal")
		if err := producer.Close(); err != nil {
			log.Fatal("error during shutdown", err)
		}
		if err := queueService.Close(); err != nil {
			log.Fatal("error during shutdown", err)
		}
		if err := db.Close(); err != nil {
			log.Fatal("error during shutdown", err)
		}
		os.Exit(0)
	}()
}

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
	db.RunMigrations(conn)

	mqConnection, err := queue.Connect()
	if err != nil {
		log.Fatal("failed connecting to queue service", err)
	}
	producer, err := queue.CreateProducer(mqConnection, "backend")
	if err != nil {
		log.Fatal("failed creating producer", err)
	}
	SetupCloseHandler(producer, mqConnection, conn)
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
