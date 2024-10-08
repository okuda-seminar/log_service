package presentation

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"log_service/internal/server/usecase"
)

type AMQPLogHandler struct {
	LogUseCase usecase.IInsertLogUseCase
}

func NewAMQPLogHandler(logUseCase usecase.IInsertLogUseCase) *AMQPLogHandler {
	return &AMQPLogHandler{
		LogUseCase: logUseCase,
	}
}

func (h *AMQPLogHandler) HandleLog(msg amqp.Delivery) {
	req, err := ParseAMQPLog(msg)
	if err != nil {
		log.Printf("Failed to parse log request: %v", err)
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
		log.Printf("Failed to insert log: %v", err)
		return
	}
	log.Printf("Received message: %s", msg.Body)
}

func ParseAMQPLog(msg amqp.Delivery) (AMQPLogRequest, error) {
	var req AMQPLogRequest

	decoder := json.NewDecoder(bytes.NewReader(msg.Body))
	err := decoder.Decode(&req)
	if err != nil {
		log.Printf("Failed to parse log request: %v", err)
		return req, err
	}
	return req, nil
}
