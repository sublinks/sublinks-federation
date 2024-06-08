package queue

import (
	"os"
	"sublinks/sublinks-federation/internal/log"
	"sublinks/sublinks-federation/internal/service"
	"sublinks/sublinks-federation/internal/worker"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue interface {
	Connect() error
	Run(serviceManager *service.ServiceManager)
	PublishMessage(queueName string, message string) error
	StartConsumer(queueData ConsumerQueue) error
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

func (q *RabbitQueue) Run(serviceManager *service.ServiceManager) {
	q.processActors(serviceManager)
	q.processObjects(serviceManager)
}

func (q *RabbitQueue) processActors(serviceManager *service.ServiceManager) {
	actorCQ := ConsumerQueue{
		QueueName: "actor_create_queue",
		Exchange:  "federation",
		RoutingKeys: map[string]worker.Worker{
			ActorRoutingKey: worker.NewActorWorker(
				q.logger,
				serviceManager.UserService(),
				serviceManager.CommunityService(),
			),
		},
	}

	err := q.StartConsumer(actorCQ)
	if err != nil {
		q.logger.Fatal("failed starting actor consumer", err)
	}
}

func (q *RabbitQueue) processObjects(serviceManager *service.ServiceManager) {
	queue := ConsumerQueue{
		QueueName: "object_create_queue",
		Exchange:  "federation",
		RoutingKeys: map[string]worker.Worker{
			PostRoutingKey: worker.NewPostWorker(
				q.logger,
				serviceManager.PostService(),
			),
			CommentRoutingKey: worker.NewCommentWorker(
				q.logger,
				serviceManager.CommentService(),
			),
		},
	}

	err := q.StartConsumer(queue)
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
