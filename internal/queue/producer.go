package queue

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	QueueName string
	*amqp.Channel
}

func (q *Queue) CreateProducer(queueName string) error {
	channelRabbitMQ, err := q.Connection.Channel()
	if err != nil {
		return err
	}
	err = q.CreateQueue(channelRabbitMQ, queueName)
	if err != nil {
		return err
	}
	q.Publishers[queueName] = &Publisher{queueName, channelRabbitMQ}
	return nil
}

func (p *Publisher) PublishMessage(message string) error {
	return p.Channel.PublishWithContext(
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
