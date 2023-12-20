package queue

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func CreateConsumer(q *amqp.Connection, queueName string) (<-chan amqp.Delivery, error) {
	channelRabbitMQ, err := q.Channel()
	if err != nil {
		return nil, err
	}
	err = CreateQueue(channelRabbitMQ, queueName)
	if err != nil {
		return nil, err
	}
	// Subscribing to QueueService1 for getting messages.
	messages, err := channelRabbitMQ.Consume(
		queueName, // queue name
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no local
		false,     // no wait
		nil,       // arguments
	)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
