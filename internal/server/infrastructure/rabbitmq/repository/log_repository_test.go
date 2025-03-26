// log_repository_test.go
package repository_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"

	"log_service/internal/server/domain"
	rabbitmqTest "log_service/internal/server/infrastructure/rabbitmq/rabbitmq_test"
	"log_service/internal/server/infrastructure/rabbitmq/repository"
)

func TestCTRSave(t *testing.T) {
	ctx := context.Background()

	rabbitContainer, connStr, err := rabbitmqTest.StartRabbitMQContainer(ctx)
	if err != nil {
		t.Fatalf("failed to start rabbitmq container: %v", err)
	}
	defer rabbitContainer.Terminate(ctx)

	conn, err := amqp.Dial(connStr)
	if err != nil {
		t.Fatalf("failed to connect to rabbitmq: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("failed to open channel: %v", err)
	}
	defer ch.Close()

	exchangeName := "test_exchange"
	routingKey := "test_key"
	queueName := "test_queue"

	err = ch.ExchangeDeclare(
		exchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to declare exchange: %v", err)
	}

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to declare queue: %v", err)
	}

	err = ch.QueueBind(
		q.Name,
		routingKey,
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to bind queue: %v", err)
	}

	repo := repository.NewLogRepository(ch, exchangeName, routingKey)

	sampleLog := &domain.CTRLog{
		EventType: "tap",
		CreatedAt: time.Now(),
		ObjectID:  "12345",
	}

	err = repo.CTRSave(ctx, sampleLog)
	assert.NoError(t, err)

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to register consumer: %v", err)
	}

	select {
	case msg := <-msgs:
		var received domain.CTRLog
		err := json.Unmarshal(msg.Body, &received)
		assert.NoError(t, err)
		assert.Equal(t, sampleLog.EventType, received.EventType)
		assert.Equal(t, sampleLog.ObjectID, received.ObjectID)
	case <-time.After(5 * time.Second):
		t.Fatal("did not receive message in time")
	}
}
