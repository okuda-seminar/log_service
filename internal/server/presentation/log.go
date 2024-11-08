package presentation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"log_service/internal/server/usecase"
	"log_service/internal/utils"
)

type AMQPLogHandler struct {
	LogUseCase usecase.IInsertLogUseCase
	Channel    *amqp.Channel
}

func NewAMQPLogHandler(logUseCase usecase.IInsertLogUseCase, ch *amqp.Channel) *AMQPLogHandler {
	return &AMQPLogHandler{
		LogUseCase: logUseCase,
		Channel:    ch,
	}
}

func (h *AMQPLogHandler) HandleLog(msg amqp.Delivery) {
	req, err := ParseAMQPLog(msg)
	if err != nil {
		h.SendResponse(utils.INVALID_ARGUMENT, fmt.Sprintf("Failed to parse log request: %v", err), msg.ReplyTo, msg.CorrelationId)
		return
	}
	logDto := &usecase.InsertLogDto{
		LogLevel:           req.LogLevel,
		Date:               req.Date,
		DestinationService: req.DestinationService,
		SourceService:      req.SourceService,
		RequestType:        req.RequestType,
		Content:            req.Content,
	}
	err = h.LogUseCase.InsertLog(context.Background(), logDto)
	if err != nil {
		h.SendResponse(utils.INTERNAL, fmt.Sprintf("Failed to insert log: %v", err), msg.ReplyTo, msg.CorrelationId)
		return
	}

	h.SendResponse(utils.OK, "OK", msg.ReplyTo, msg.CorrelationId)
}

func ParseAMQPLog(msg amqp.Delivery) (AMQPLogRequest, error) {
	var req AMQPLogRequest

	decoder := json.NewDecoder(bytes.NewReader(msg.Body))
	err := decoder.Decode(&req)
	if err != nil {
		return req, err
	}
	return req, nil
}

func (h *AMQPLogHandler) SendResponse(statusCode int, message, key, corrID string) {
	res := &AmqpLogResponse{
		StatusCode: statusCode,
		Message:    message,
	}
	bytes, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	err = h.Channel.Publish(
		"",    // exchange
		key,   // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrID,
			Body:          bytes,
		})
	if err != nil {
		panic(err)
	}
}
