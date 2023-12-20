package queue

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

func CreateProducer(q *amqp.Connection, queueName string) (*amqp.Channel, error) {
	channelRabbitMQ, err := q.Channel()
	if err != nil {
		return nil, err
	}
	err = CreateQueue(channelRabbitMQ, queueName)
	if err != nil {
		return nil, err
	}
	return channelRabbitMQ, nil
}

func PublishMessage(q *amqp.Channel, message string) error {
	return q.PublishWithContext(
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
