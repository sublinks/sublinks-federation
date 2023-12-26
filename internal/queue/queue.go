package queue

import (
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Connect() (*amqp.Connection, error) {
	// Get the connection string from the environment variable
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")
	// Create a new RabbitMQ connection.
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		return nil, err
	}
	return connectRabbitMQ, nil
}

func CreateQueue(channelRabbitMQ *amqp.Channel, queueName string) error {
	// With the instance and declare Queues that we can
	// publish and subscribe to.
	_, err := channelRabbitMQ.QueueDeclare(
		queueName, // queue name
		true,      // durable
		false,     // auto delete
		false,     // exclusive
		false,     // no wait
		nil,       // arguments
	)
	return err
}
