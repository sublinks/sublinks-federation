package queue

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
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

func (q *RabbitQueue) StartConsumer(queueData ConsumerQueue, ctx context.Context) error {
	err := q.createConsumer(queueData)
	if err != nil {
		return err
	}
	messages, ok := q.consumers[queueData.QueueName]
	if !ok {
		return fmt.Errorf("consumer not found")
	}

	errGroup, ctx := errgroup.WithContext(ctx)
	for {
		select {
		case <-ctx.Done():
			q.logger.Debug("consumer context canceled")
			return errGroup.Wait()
		case message, ok := <-messages:
			if !ok {
				q.logger.Error("consumer channel closed", errors.New("consumer channel closed"))
				return errGroup.Wait()
			}
			msg := message
			errGroup.Go(func() error {
				cbWorker, ok := queueData.RoutingKeys[msg.RoutingKey]
				q.logger.Info(fmt.Sprintf("consumer got message from routing key: %s", msg.RoutingKey))
				if !ok {
					return fmt.Errorf("%s not implemented as valid routing key", msg.RoutingKey)
				}

				err := cbWorker.Process(msg.Body)
				if err != nil {
					nackErr := msg.Acknowledger.Nack(msg.DeliveryTag, false, true)
					if nackErr != nil {
						return fmt.Errorf("error nack'ing the message: %s", nackErr.Error())
					}
					return fmt.Errorf("error processing message body: %s", err.Error())
				}

				ackErr := msg.Acknowledger.Ack(msg.DeliveryTag, false)
				if ackErr != nil {
					return fmt.Errorf("error ack'ing the message: %s", ackErr.Error())
				}
				return nil
			})
		}
	}
}
