package queue

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type publisher struct {
	QueueName string
	*amqp.Channel
}

func (q *RabbitQueue) createProducer(queueName string) error {
	channelRabbitMQ, err := q.Connection.Channel()
	if err != nil {
		return err
	}
	err = q.createQueue(channelRabbitMQ, queueName)
	if err != nil {
		return err
	}
	q.publishers[queueName] = &publisher{queueName, channelRabbitMQ}
	return nil
}

func (q *RabbitQueue) PublishMessage(queueName string, message string) error {
	publisher, ok := q.publishers[queueName]
	if !ok {
		_ = q.createProducer(queueName)
		publisher = q.publishers[queueName]
	}
	return publisher.Channel.PublishWithContext(
		context.TODO(),
		"backend", // exchange
		"",        // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}
