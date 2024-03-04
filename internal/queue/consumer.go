package queue

import (
	"errors"
)

func (q *RabbitQueue) createConsumer(queueName string, exchangeName string) error {
	channelRabbitMQ, err := q.Connection.Channel()
	if err != nil {
		return err
	}

	err = channelRabbitMQ.ExchangeDeclare(
		exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for _, key := range q.routingKeys {
		err = q.createQueue(channelRabbitMQ, queueName)
		if err != nil {
			return err
		}

		err = channelRabbitMQ.QueueBind(queueName, key, exchangeName, false, nil)
		if err != nil {
			return err
		}
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
		return err
	}

	q.consumers[queueName] = messages

	return nil
}

// TODO: Implement a way to either pass a callback function or return messages/chan
func (q *RabbitQueue) StartConsumer(queueName string, exchangeName string) error {
	err := q.createConsumer(queueName, exchangeName)
	if err != nil {
		return err
	}

	messages, ok := q.consumers[queueName]
	if !ok {
		return errors.New("consumer not found")
	}

	go func() {
		for message := range messages {
			q.logger.Printf(" > Received message: %s\n", message.Body)
		}
	}()

	return nil
}
