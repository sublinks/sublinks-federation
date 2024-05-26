package queue

import (
	"fmt"
	"sublinks/sublinks-federation/internal/worker"

	"golang.org/x/sync/errgroup"
)

type ConsumerQueue struct {
	Exchange    string
	QueueName   string
	RoutingKeys map[string]worker.Worker
}

func (q *RabbitQueue) createConsumer(queueData ConsumerQueue) error {
	channelRabbitMQ, err := q.Connection.Channel()
	if err != nil {
		return err
	}
	err = q.createQueue(channelRabbitMQ, queueData.QueueName)
	if err != nil {
		return err
	}

	for routingKey := range queueData.RoutingKeys {
		err = channelRabbitMQ.QueueBind(
			queueData.QueueName, // queue name
			routingKey,          // routing key
			queueData.Exchange,  // exchange
			false,
			nil)
		if err != nil {
			return err
		}
	}

	// Subscribing to QueueService1 for getting messages.
	messages, err := channelRabbitMQ.Consume(
		queueData.QueueName, // queue name
		"",                  // consumer
		false,               // auto-ack
		false,               // exclusive
		false,               // no local
		false,               // no wait
		nil,                 // arguments
	)
	if err != nil {
		return err
	}
	q.consumers[queueData.QueueName] = messages
	return nil
}

func (q *RabbitQueue) StartConsumer(queueData ConsumerQueue) error {
	err := q.createConsumer(queueData)
	if err != nil {
		return err
	}
	messages, ok := q.consumers[queueData.QueueName]
	if !ok {
		return fmt.Errorf("consumer not found")
	}

	errGroup := new(errgroup.Group)
	for message := range messages {
		errGroup.Go(func() error {
			cbWorker, ok := queueData.RoutingKeys[message.RoutingKey]
			if !ok {
				return fmt.Errorf("%s not implemented as valid routing key", message.RoutingKey)
			}

			err := cbWorker.Process(message.Body)

			if err != nil {
				err = message.Acknowledger.Nack(message.DeliveryTag, false, true)
				if err != nil {
					return fmt.Errorf("error nack'ing the message: %s", err.Error())
				}
				return fmt.Errorf("error processing message body: %s", err.Error())
			}

			err = message.Acknowledger.Ack(message.DeliveryTag, false)
			if err != nil {
				return fmt.Errorf("error ack'ing the message: %s", err.Error())
			}
			return nil
		})
	}

	if err := errGroup.Wait(); err != nil {
		return err
	}

	return nil
}
