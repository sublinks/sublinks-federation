package main

import (
	"fmt"
	"os"
	"sublinks/sublinks-federation/internal/db"
	"sublinks/sublinks-federation/internal/http"
	"sublinks/sublinks-federation/internal/log"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
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
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")
	// Create a new RabbitMQ connection.
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		panic(err)
	}
	defer connectRabbitMQ.Close()
	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	defer channelRabbitMQ.Close()
	// With the instance and declare Queues that we can
	// publish and subscribe to.
	_, err = channelRabbitMQ.QueueDeclare(
		"QueueService1", // queue name
		true,            // durable
		false,           // auto delete
		false,           // exclusive
		false,           // no wait
		nil,             // arguments
	)
	if err != nil {
		panic(err)
	}
	http.RunServer()
	os.Exit(0)
}
