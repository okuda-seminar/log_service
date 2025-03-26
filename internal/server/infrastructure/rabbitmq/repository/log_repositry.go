// log_repository.go
package repository

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"

	"log_service/internal/server/domain"
)

// LogRepository is a struct that represents the repository for logging.
// It contains a channel, exchange, and routing key.
type LogRepository struct {
	channel    *amqp.Channel
	exchange   string
	routingKey string
}

// NewLogRepository generates a new LogRepository instance.
// It takes a channel, exchange, and routing key as arguments.
func NewLogRepository(ch *amqp.Channel, exchange, routingKey string) *LogRepository {
	return &LogRepository{
		channel:    ch,
		exchange:   exchange,
		routingKey: routingKey,
	}
}

// CTRSave save CTRLog into RabbitMQ.
// It takes a context and a CTRLog object from the domain package as arguments.
func (r *LogRepository) CTRSave(ctx context.Context, ctrLog *domain.CTRLog) error {
	body, err := json.Marshal(ctrLog)
	if err != nil {
		return err
	}

	return r.channel.Publish(
		r.exchange,   // exchange name
		r.routingKey, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
