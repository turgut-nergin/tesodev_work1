package models

import (
	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Channel *amqp091.Channel
	Queue   amqp091.Queue
	Message string
}
