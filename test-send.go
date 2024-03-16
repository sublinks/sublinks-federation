package main

import (
	"context"
	"fmt"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"sublinks/sublinks-federation/internal/log"

	"github.com/joho/godotenv"
)

func failOnError(err error, msg string) {
	if err != nil {
		logger := log.NewLogger("test-send")
		logger.Panic().Err(err).Msg(msg)
	}
}

func main() {
	// bootstrap logger
	logger := log.NewLogger("test-send")

	// Load connection string from .env file
	err := godotenv.Load()
	if err != nil {
		logger.Warn(fmt.Sprintf("failed to load env, %v", err))
	}

	// Get the connection string from the environment variable
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")
	// Create a new RabbitMQ connection.
	conn, err := amqp.Dial(amqpServerURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sendMessage(conn, ctx, "actor.create", `
{
  "actor_type": "Person",
  "id": "https://sublinks.org/u/lazyguru",
  "username": "lazyguru",
  "name": "Lazyguru",
  "bio": "I am a lazy guru",
  "matrix_user_id": "@lazyguru:discuss.online",
  "private_key": "some super secure string",
  "public_key": "the public key"
}`)

	sendMessage(conn, ctx, "actor.create", `
{
  "actor_type": "Group",
  "id": "https://sublinks.org/c/test_community",
  "username": "test_community",
  "name": "Test Community",
  "private_key": "some super secure string",
  "public_key": "the public key",
  "sensitive": false
}`)

	sendMessage(conn, ctx, "post.create", `
{
  "id": "https://sublinks.org/post/test-post-1",
  "title": "This is a post",
  "content": "I am a lazy guru",
  "published": "2021-08-01T12:34:59Z",
  "community": "https://sublinks.org/c/test_community",
  "author": "https://sublinks.org/u/lazyguru"
}`)
}

func sendMessage(conn *amqp.Connection, ctx context.Context, routingKey string, body string) {
	// bootstrap logger
	logger := log.NewLogger("test-send")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	err = ch.PublishWithContext(ctx,
		"federation", // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Timestamp:   time.Now(),
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	logger.Info(fmt.Sprintf("[x] Sent %s", body))

	ch2, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch2.Close()

	body = `
{
  "id": "test-post-1",
  "title": "This is a post",
  "content": "I am a lazy guru",
  "published": "2021-08-01T12:34:59Z",
  "community": "https://sublinks.org/c/test-community",
  "author": "https://sublinks.org/u/lazyguru"
}`

	err = ch2.PublishWithContext(ctx,
		"federation",  // exchange
		"post.create", // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Timestamp:   time.Now(),
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	logger.Info(fmt.Sprintf("[x] Sent %s", body))
}
