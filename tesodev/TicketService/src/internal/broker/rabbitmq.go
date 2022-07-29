package broker

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/turgut-nergin/tesodev_work1/config"
	"github.com/turgut-nergin/tesodev_work1/internal/errors"
	"github.com/turgut-nergin/tesodev_work1/internal/repository"
)

func CreateConnection(config config.RabbitMQConfig) *amqp.Connection {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", config.UserName, config.Password, config.Host, config.Port)
	conn, err := amqp.Dial(url)
	errors.FailOnError(err, "Failed to connect to RabbitMQ")
	return conn
}

func CreateChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	errors.FailOnError(err, "Failed to open a channel")
	return ch
}

func CreateQueue(ch *amqp.Channel) amqp.Queue {
	q, err := ch.QueueDeclare(
		"signal", // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	errors.FailOnError(err, "Failed to declare a queue")
	return q
}

func RegisterConsumer(ch *amqp.Channel, q amqp.Queue) <-chan amqp.Delivery {
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	errors.FailOnError(err, "Failed to register a consumer")
	return msgs
}

func HandleDelivery(msgs <-chan amqp.Delivery, repository repository.Repositories) {
	go func() {
		for d := range msgs {
			repository.TicketRepository.UpdateTicket(string(d.Body))
		}
	}()
}
