package queue

import (
	"os"

	"sublinks/sublinks-federation/internal/log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue interface {
	Connect() error
	PublishMessage(queueName string, message string) error
	StartConsumer(queueName string) error
	Close()
}

type RabbitQueue struct {
	*amqp.Connection
	publishers  map[string]*publisher
	consumers   map[string]<-chan amqp.Delivery
	routingKeys []string
	logger      *log.Log
}

func NewQueue(logger *log.Log, routingKeys []string) Queue {
	return &RabbitQueue{
		logger:      logger,
		publishers:  make(map[string]*publisher),
		consumers:   make(map[string]<-chan amqp.Delivery),
		routingKeys: routingKeys,
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

func (q *RabbitQueue) createQueue(channelRabbitMQ *amqp.Channel, queueName string) error {
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
	for _, publisher := range q.publishers {
		publisher.Close()
	}
	q.Connection.Close()
}
