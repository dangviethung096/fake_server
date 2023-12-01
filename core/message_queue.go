package core

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type messageQueue struct {
	connection *amqp.Connection
}

type MessageQueueSession struct {
	channel    *amqp.Channel
	connection *messageQueue
	config     QueueConfig
}

type QueueConfig struct {
	ExchangeName string
	QueueName    string
	RouteKey     string
	Kind         string
	Durable      bool
	AutoDelete   bool
	Exclusive    bool
	NoWait       bool
	Args         amqp.Table
	ConsumerTag  string
	AutoAck      bool
}

func connectRabbitMQ() *messageQueue {
	fmt.Printf("Connect to RabbitMQ: %s\n", Config.RabbitMQ.AMQPServerURL)
	// Connect to RabbitMQ
	conn, err := amqp.Dial(Config.RabbitMQ.AMQPServerURL)
	if err != nil {
		log.Fatalf("Could not establish connection with RabbitMQ: %s", err.Error())
	}

	connection := &messageQueue{conn}
	// Retry to reconnect if
	go func(con *messageQueue) {
		for err := range conn.NotifyClose(make(chan *amqp.Error)) {
			LoggerInstance.Error("Connection to rabbitmq is disconnected: retry to connect %s", err.Error())
			con.retryConnect()
		}
	}(connection)

	return connection
}

func (mq *messageQueue) retryConnect() Error {
	// Close old connection
	mq.connection.Close()
	// Open new connection
	connection, err := amqp.Dial(Config.RabbitMQ.AMQPServerURL)
	if err != nil {
		LoggerInstance.Error("Could not establish connection with RabbitMQ: %s", err.Error())
		return ERROR_CANNOT_CONNECT_RABBITMQ
	}

	mq.connection = connection
	return nil
}

func (mq *messageQueue) CreateSession(config QueueConfig) (*MessageQueueSession, Error) {
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
		LoggerInstance.Error("Error when declare queue: %s", originalErr.Error())
		return nil, ERROR_CANNOT_DECLARE_QUEUE
	}

	originalErr = channel.ExchangeDeclare(
		config.ExchangeName,
		config.Kind,
		config.Durable,
		config.AutoDelete,
		config.Exclusive,
		config.NoWait,
		config.Args,
	)

	if originalErr != nil {
		LoggerInstance.Error("Error when declare exchange: %s", originalErr.Error())
		return nil, ERROR_CANNOT_DECLARE_EXCHANGE
	}

	originalErr = channel.QueueBind(
		config.QueueName,
		config.RouteKey,
		config.ExchangeName,
		config.NoWait,
		config.Args,
	)

	if originalErr != nil {
		LoggerInstance.Error("Error when bind queue: %s", originalErr.Error())
		return nil, ERROR_CANNOT_BIND_QUEUE
	}

	session := &MessageQueueSession{
		channel:    channel,
		connection: mq,
		config:     config,
	}

	go func(sess *MessageQueueSession) {
		for err := range sess.channel.NotifyClose(make(chan *amqp.Error)) {
			LoggerInstance.Error("Error: %s", err.Error())
			// Close old connection
			sess.connection.retryConnect()
			// Close session
			sess.channel.Close()
			// Retry create a new session
			for !sess.recreateSession() {
				time.Sleep(time.Second * time.Duration(Config.RabbitMQ.RetryTime))
			}
		}
	}(session)

	return session, nil
}

/*
* recreateSession
 */
func (session *MessageQueueSession) recreateSession() bool {
	channel, err := session.connection.connection.Channel()
	if err != nil {
		LoggerInstance.Error("Could not open channel with RabbitMQ: %s", err.Error())
		return false
	}

	_, originalErr := channel.QueueDeclare(
		session.config.QueueName,
		session.config.Durable,
		session.config.AutoDelete,
		session.config.Exclusive,
		session.config.NoWait,
		session.config.Args,
	)

	if originalErr != nil {
		LoggerInstance.Error("Error when declare queue: %s", originalErr.Error())
		return false
	}

	originalErr = channel.ExchangeDeclare(
		session.config.ExchangeName,
		session.config.Kind,
		session.config.Durable,
		session.config.AutoDelete,
		session.config.Exclusive,
		session.config.NoWait,
		session.config.Args,
	)

	if originalErr != nil {
		LoggerInstance.Error("Error when declare exchange: %s", originalErr.Error())
		return false
	}

	originalErr = channel.QueueBind(
		session.config.QueueName,
		session.config.RouteKey,
		session.config.ExchangeName,
		session.config.NoWait,
		session.config.Args,
	)

	if originalErr != nil {
		LoggerInstance.Error("Error when bind queue: %s", originalErr.Error())
		return false
	}

	session.channel = channel

	return true
}

func (mqs *MessageQueueSession) CloseSession() {
	if mqs.channel != nil {
		mqs.channel.Close()
	}
}

func (mqs *MessageQueueSession) Publish(body []byte) Error {
	err := mqs.channel.PublishWithContext(
		coreContext,
		mqs.config.ExchangeName,
		mqs.config.RouteKey,
		false,
		false,
		amqp.Publishing{
			ContentType: CONTENT_TYPE_TEXT,
			Body:        body,
		},
	)

	if err != nil {
		LoggerInstance.Error("Publish error: %v", err)
		return ERROR_SERVER_ERROR
	}

	return nil
}

type RabbitmqMessage struct {
	Body []byte
}

type ConsumerHandler func(msg RabbitmqMessage)

func (mqs *MessageQueueSession) Consume(handler ConsumerHandler) Error {
	consumerTag := mqs.config.ConsumerTag
	if consumerTag == BLANK {
		consumerTag = DEFAULT_CONSUMER_TAG
	}
	messages, err := mqs.channel.Consume(mqs.config.QueueName, consumerTag, true, false, false, false, nil)
	if err != nil {
		LoggerInstance.Error("Error when consume messages: %s", err.Error())
		return ERROR_CANNOT_CONSUME_MESSAGES_FROM_RABBITMQ
	}

	// Handle message from rabbitmq
	go func(c <-chan amqp.Delivery) {
		for message := range c {
			handler(RabbitmqMessage{
				Body: message.Body,
			})
		}
	}(messages)

	return nil
}
