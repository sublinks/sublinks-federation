package queue

import (
	"os"

	"sublinks/sublinks-federation/internal/log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue struct {
	*amqp.Connection
	Publishers map[string]*Publisher
	Consumers  map[string]<-chan amqp.Delivery
	Logger     *log.Log
}

func NewQueue(logger *log.Log) *Queue {
	return &Queue{Logger: logger}
}

func (q *Queue) Connect() error {
	// Get the connection string from the environment variable
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")
	// Create a new RabbitMQ connection.
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		return err
	}
	q.Connection = connectRabbitMQ
	return nil
}

func (q *Queue) CreateQueue(channelRabbitMQ *amqp.Channel, queueName string) error {
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

func (q *Queue) Close() {
	for _, publisher := range q.Publishers {
		publisher.Close()
	}
	q.Connection.Close()
}
