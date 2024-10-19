package rabbitmq

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"

	"log_service/internal/server/infrastructure/rabbitmq"
)

func Connect() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	return conn, ch, nil
}

func Publish(ch *amqp.Channel, input []byte) error {
	err := ch.Publish(
		"",                  // exchange
		rabbitmq.QUEUE_NAME, // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(input),
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
		return err
	}
	return nil
}
