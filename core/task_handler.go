package core

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type TaskInfo struct {
	Data []byte
}

type TaskHandler func(ctx *Context, task TaskInfo)

/*
* Handle task: handle a task message is published from message queue
 */
func HandleTask(taskQueueName string, handler TaskHandler) Error {
	queueConfig := QueueConfig{
		ExchangeName: BLANK,
		QueueName:    fmt.Sprintf("%s%s", TASK_PREFIX_QUEUE_NAME, taskQueueName),
		RouteKey:     BLANK,
		Kind:         MESSAGE_QUEUE_KIND_DIRECT,
		AutoAck:      true,
		Durable:      false,
		AutoDelete:   false,
		Exclusive:    false,
		NoWait:       false,
		Args:         nil,
	}

	session, err := MessageQueue().CreateSimpleSession(queueConfig)

	if err != nil {
		LoggerInstance.Error("Create message queue session fail: %s", err.Error())
		return err
	}

	messages, errConsume := session.channel.Consume(session.config.QueueName, BLANK, true, session.config.Exclusive, false, session.config.NoWait, nil)
	if errConsume != nil {
		LoggerInstance.Error("Error when handle task: %s", err.Error())
		return ERROR_CANNOT_CONSUME_MESSAGES_FROM_RABBITMQ
	}

	go func(messages <-chan amqp.Delivery) {
		for message := range messages {
			ctx := GetContextWithTimeout(Config.GetTaskTimeout())
			ctx.LogInfo("Start handle task: %s", queueConfig.QueueName)
			handler(ctx, TaskInfo{
				Data: message.Body,
			})
			ctx.LogInfo("End handle task: %s", queueConfig.QueueName)
		}
	}(messages)

	return nil
}
