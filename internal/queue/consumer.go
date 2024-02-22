package queue

import (
	"errors"
	"sublinks/sublinks-federation/internal/worker"
)

func (q *RabbitQueue) createConsumer(queueData ConsumerQueue) error {
	channelRabbitMQ, err := q.Connection.Channel()
	if err != nil {
		return err
	}
	err = q.createQueue(channelRabbitMQ, queueData.QueueName)
	if err != nil {
		return err
	}
	err = channelRabbitMQ.QueueBind(
		queueData.QueueName,  // queue name
		queueData.RoutingKey, // routing key
		queueData.Exchange,   // exchange
		false,
		nil)
	if err != nil {
		return err
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

type ConsumerQueue struct {
	Exchange   string
	QueueName  string
	RoutingKey string
}

// TODO: Implement a way to either pass a callback function or return messages/chan
func (q *RabbitQueue) StartConsumer(queueData ConsumerQueue, worker worker.Worker) error {
	err := q.createConsumer(queueData)
	if err != nil {
		return err
	}
	messages, ok := q.consumers[queueData.QueueName]
	if !ok {
		return errors.New("consumer not found")
	}
	go func() {
		for message := range messages {
			err := worker.Process(message.Body)
			if err != nil {
				err = message.Acknowledger.Nack(message.DeliveryTag, false, true)
				if err != nil {
					panic(err) // I know this isn't good. Will need to fix it
				}
				continue
			}
			err = message.Acknowledger.Ack(message.DeliveryTag, false)
			if err != nil {
				panic(err) // I know this isn't good. Will need to fix it
			}
		}
	}()
	return nil
}
