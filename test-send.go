package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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
	log.Printf(" [x] Sent %s\n", body)
}
