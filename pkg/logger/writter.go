package logger

import (
	"context"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"

	"log_service/internal/client/presentation"
	serverPresentation "log_service/internal/server/presentation"
)

type Writter struct {
	ch              *amqp.Channel
	logPresentation presentation.ILogPresentation
	queue           string
	msgs            <-chan amqp.Delivery
}

func NewWritter(ctx context.Context, ch *amqp.Channel) (*Writter, error) {
	logPresentation := presentation.NewLogPresentation(ch)
	msgs, queue, err := logPresentation.Consume()
	if err != nil {
		return nil, err
	}
	w := &Writter{
		ch:              ch,
		logPresentation: logPresentation,
		queue:           queue,
		msgs:            msgs,
	}
	return w, nil
}

func (w *Writter) Write(p []byte) (n int, err error) {
	req := serverPresentation.AMQPLogRequest{
		Date:    time.Now(),
		Content: string(p),
	}

	id := uuid.New().String()
	if err := w.logPresentation.Publish(context.Background(), w.queue, id, req); err != nil {
		return 0, err
	}
	if err := w.logPresentation.Serve(w.msgs, id); err != nil {
		return 0, err
	}
	return len(p), nil
}
