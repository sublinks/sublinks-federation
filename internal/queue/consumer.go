package queue

import (
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"sublinks/sublinks-federation/internal/worker"
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

	for routingKey, _ := range queueData.RoutingKeys {
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
		return errors.New("consumer not found")
	}

	errGroup := new(errgroup.Group)
	for message := range messages {
		errGroup.Go(func() error {
			cbWorker, ok := queueData.RoutingKeys[message.RoutingKey]
			if !ok {
				return errors.New(fmt.Sprintf("%s not implemented as valid routing key", message.RoutingKey))
			}

			err := cbWorker.Process(message.Body)

			if err != nil {
				err = message.Acknowledger.Nack(message.DeliveryTag, false, true)
				if err != nil {
					return errors.New(fmt.Sprintf("error nack'ing the message: %s", err.Error()))
				}
				return errors.New(fmt.Sprintf("error processing message body: %s", err.Error()))
			}

			err = message.Acknowledger.Ack(message.DeliveryTag, false)
			if err != nil {
				return errors.New(fmt.Sprintf("error ack'ing the message: %s", err.Error()))
			}
			return nil
		})
	}

	if err := errGroup.Wait(); err != nil {
		return err
	}

	return nil
}
