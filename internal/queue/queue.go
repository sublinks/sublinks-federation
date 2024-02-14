package queue

import (
	"os"

	"sublinks/sublinks-federation/internal/log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue interface {
	Connect() error
	CreateProducer(queueName string) error
	CreateConsumer(queueName string) error
	StartConsumer(queueName string)
	Close()
}

type RabbitQueue struct {
	*amqp.Connection
	Publishers map[string]MessagePublisher
	Consumers  map[string]<-chan amqp.Delivery
	Logger     *log.Log
}

func NewQueue(logger *log.Log) Queue {
	return &RabbitQueue{
		Logger:     logger,
		Publishers: make(map[string]MessagePublisher),
		Consumers:  make(map[string]<-chan amqp.Delivery),
	}
}

func (q *RabbitQueue) Connect() error {
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

func (q *RabbitQueue) CreateQueue(channelRabbitMQ *amqp.Channel, queueName string) error {
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

func (q *RabbitQueue) Close() {
	for _, publisher := range q.Publishers {
		publisher.Close()
	}
	q.Connection.Close()
}
