package core

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleMessageQueueSession struct {
	channel    *amqp.Channel
	connection *messageQueue
	config     QueueConfig
}

func (mq *messageQueue) CreateSimpleSession(config QueueConfig) (*SimpleMessageQueueSession, Error) {
	channel, err := mq.connection.Channel()
	if err != nil {
		LoggerInstance.Error("Could not open channel with RabbitMQ: %s", err.Error())
		return nil, ERROR_CANNOT_CREATE_RABBITMQ_CHANNEL
	}

	_, originalErr := channel.QueueDeclare(
		config.QueueName,
		config.Durable,
		config.AutoDelete,
		config.Exclusive,
		config.NoWait,
		config.Args,
	)

	if originalErr != nil {
		LoggerInstance.Error("Error when declare queue: %s", err.Error())
		return nil, ERROR_CANNOT_DECLARE_QUEUE
	}

	session := &SimpleMessageQueueSession{
		channel:    channel,
		connection: mq,
		config:     config,
	}

	// Retry to connect
	go func(sess *SimpleMessageQueueSession) {
		for err := range sess.channel.NotifyClose(make(chan *amqp.Error)) {
			LoggerInstance.Error("Channel disconnected: retry to connect: %s", err.Error())
			for !sess.recreateSession() {
				time.Sleep(time.Second * time.Duration(Config.RabbitMQ.RetryTime))
			}

		}
	}(session)

	return session, nil
}

func (mqs *SimpleMessageQueueSession) CloseSession() {
	if mqs.channel != nil {
		mqs.channel.Close()
	}
}

func (mqs *SimpleMessageQueueSession) recreateSession() bool {
	channel, err := mqs.connection.connection.Channel()
	if err != nil {
		LoggerInstance.Error("Could not open channel with RabbitMQ: %s", err.Error())
		return false
	}

	_, originalErr := channel.QueueDeclare(
		mqs.config.QueueName,
		mqs.config.Durable,
		mqs.config.AutoDelete,
		mqs.config.Exclusive,
		mqs.config.NoWait,
		mqs.config.Args,
	)

	if originalErr != nil {
		LoggerInstance.Error("Error when declare queue: %s", err.Error())
		return false
	}

	mqs.channel = channel
	return true
}

func (mqs *SimpleMessageQueueSession) Publish(body []byte) Error {
	err := mqs.channel.PublishWithContext(
		coreContext,
		mqs.config.ExchangeName,
		mqs.config.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: CONTENT_TYPE_TEXT,
			Body:        body,
		},
	)

	if err != nil {
		LoggerInstance.Error("Publish message: %v", err)
		return ERROR_SERVER_ERROR
	}

	return nil
}

func (mqs *SimpleMessageQueueSession) Consume(handler func(msg RabbitmqMessage)) Error {
	messages, err := mqs.channel.Consume(mqs.config.QueueName, BLANK, true, mqs.config.Exclusive, false, mqs.config.NoWait, nil)
	if err != nil {
		LoggerInstance.Error("Error when consume messages: %s", err.Error())
		return ERROR_CANNOT_CONSUME_MESSAGES_FROM_RABBITMQ
	}

	go func(messages <-chan amqp.Delivery) {
		for message := range messages {
			handler(RabbitmqMessage{
				Body: message.Body,
			})
		}
	}(messages)

	return nil
}
