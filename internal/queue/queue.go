package queue

import (
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"sublinks/sublinks-federation/internal/db"
	"sublinks/sublinks-federation/internal/log"
)

type Queue interface {
	Connect() error
	Run(conn db.Database)
	PublishMessage(queueName string, message string) error
	StartConsumer(queueData ConsumerQueue, conn db.Database) error
	Status() map[string]map[string]bool
	Close()
}

type RabbitQueue struct {
	*amqp.Connection
	publishers map[string]*publisher
	consumers  map[string]<-chan amqp.Delivery
	logger     *log.Log
}

func NewQueue(logger *log.Log) Queue {
	return &RabbitQueue{
		logger:     logger,
		publishers: make(map[string]*publisher),
		consumers:  make(map[string]<-chan amqp.Delivery),
	}
}

func (q *RabbitQueue) Status() map[string]map[string]bool {
	status := make(map[string]map[string]bool)
	publisherStatus := make(map[string]bool)
	consumerStatus := make(map[string]bool)
	for publisherName, publisher := range q.publishers {
		publisherStatus[publisherName] = !publisher.IsClosed()
	}
	status["publishers"] = publisherStatus
	for consumerName, consumer := range q.consumers {
		consumerStatus[consumerName] = consumer != nil
	}
	status["consumers"] = consumerStatus
	return status
}

func (q *RabbitQueue) Run(conn db.Database) {
	q.processActors(conn)
	q.processObjects(conn)
}

func (q *RabbitQueue) processActors(conn db.Database) {
	actorCQ := ConsumerQueue{
		QueueName:   "actor_create_queue",
		Exchange:    "federation",
		RoutingKeys: []string{"actor.create"},
	}

	err := q.StartConsumer(actorCQ, conn)
	if err != nil {
		q.logger.Fatal("failed starting actor consumer", err)
	}
}

func (q *RabbitQueue) processObjects(conn db.Database) {
	queue := ConsumerQueue{
		QueueName:   "object_create_queue",
		Exchange:    "federation",
		RoutingKeys: []string{"post.create", "comment.create"},
	}

	err := q.StartConsumer(queue, conn)
	if err != nil {
		q.logger.Fatal("failed starting object consumer", err)
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
