package queue

import (
	"os"

	"sublinks/sublinks-federation/internal/db"
	"sublinks/sublinks-federation/internal/log"
	"sublinks/sublinks-federation/internal/repository"
	"sublinks/sublinks-federation/internal/worker"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue interface {
	Connect() error
	Run(conn db.Database)
	PublishMessage(queueName string, message string) error
	StartConsumer(queueData ConsumerQueue, worker worker.Worker) error
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

func (q *RabbitQueue) Run(conn db.Database) {
	q.processActors(conn)
	q.processPosts(conn)
}

func (q *RabbitQueue) processActors(conn db.Database) {
	actorCQ := ConsumerQueue{
		QueueName:  "actor_create_queue",
		Exchange:   "federation",
		RoutingKey: "actor.create",
	}

	aw := worker.ActorWorker{
		Logger:     q.logger,
		Repository: repository.NewRepository(conn),
	}

	err := q.StartConsumer(actorCQ, &aw)
	if err != nil {
		q.logger.Fatal("failed starting actor consumer", err)
	}
}

func (q *RabbitQueue) processPosts(conn db.Database) {
	postCQ := ConsumerQueue{
		QueueName:  "post_queue",
		Exchange:   "federation",
		RoutingKey: "post.create",
	}

	aw := worker.PostWorker{
		Logger:     q.logger,
		Repository: repository.NewRepository(conn),
	}

	err := q.StartConsumer(postCQ, &aw)
	if err != nil {
		q.logger.Fatal("failed starting post consumer", err)
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
