package queue

func (q *Queue) CreateConsumer(queueName string) error {
	channelRabbitMQ, err := q.Connection.Channel()
	if err != nil {
		return err
	}
	err = q.CreateQueue(channelRabbitMQ, queueName)
	if err != nil {
		return err
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
	q.Consumers[queueName] = messages
	return nil
}

func (q *Queue) StartConsumer(queueName string) {
	messages, ok := q.Consumers[queueName]
	if !ok {
		return
	}
	go func() {
		for message := range messages {
			q.Logger.Printf(" > Received message: %s\n", message.Body)
		}
	}()
}
