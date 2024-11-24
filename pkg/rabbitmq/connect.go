package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"

	"log_service/internal/client/infrastructure/rabbitmq"
)

func Connect() (*amqp.Connection, *amqp.Channel, error) {
	conn, ch, err := rabbitmq.Connect()
	return conn, ch, err
}
