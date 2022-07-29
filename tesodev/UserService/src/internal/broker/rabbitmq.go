package broker

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/turgut-nergin/tesodev_work1/config"
	"github.com/turgut-nergin/tesodev_work1/internal/errors"
	"github.com/turgut-nergin/tesodev_work1/internal/models"
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

func PublishConsumer(rabbitMQ *models.RabbitMQ) error {
	err := rabbitMQ.Channel.Publish(
		"",                  // exchange
		rabbitMQ.Queue.Name, // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(rabbitMQ.Message),
		})
	return err
}

// func SendMessage(message string) {

// 	config := config.RabbitMQEnvConfig["local"]
// 	conn := createConnection(config)
// 	defer conn.Close()

// 	ch := createChannel(conn)
// 	defer ch.Close()

// 	q := createQueue(ch)

// 	err := publishConsumer(ch, q, message)
// 	errors.FailOnError(err, "Failed to publish a message")
// }
