package presentation

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"log_service/internal/server/infrastructure/rabbitmq"
	"log_service/internal/server/presentation"
)

type ILogPresentation interface {
	Publish(ctx context.Context, queueName string, req presentation.AMQPLogRequest) error
	Consume() (<-chan amqp.Delivery, string, error)
	Serve(msgs <-chan amqp.Delivery) error
}

type LogPresentation struct {
	ch *amqp.Channel
	id string
}

func NewLogPresentation(ch *amqp.Channel, id string) *LogPresentation {
	return &LogPresentation{
		ch: ch,
		id: id,
	}
}

func (r *LogPresentation) Publish(ctx context.Context, queueName string, req presentation.AMQPLogRequest) error {
	bytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	err = r.ch.Publish(
		"",                  // exchange
		rabbitmq.QUEUE_NAME, // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			ReplyTo:       queueName,
			CorrelationId: r.id,
			Body:          bytes,
		})
	if err != nil {
		return err
	}
	return nil
}

func (r *LogPresentation) Consume() (<-chan amqp.Delivery, string, error) {
	q, err := r.ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, "", err
	}
	msgs, err := r.ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, "", err
	}
	return msgs, q.Name, nil
}

func (r *LogPresentation) Serve(msgs <-chan amqp.Delivery) error {
	for d := range msgs {
		if d.CorrelationId == r.id {
			res := &presentation.AmqpLogResponse{}
			if err := json.Unmarshal(d.Body, res); err != nil {
				return err
			}
			log.Printf(
				"Received response: code: %d, message: %v",
				res.StatusCode,
				res.Message,
			)
			break
		}
	}
	return nil
}
