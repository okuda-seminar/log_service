package presentation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"

	"log_service/internal/server/usecase"
	"log_service/internal/utils"
)

type AMQPLogHandler struct {
	LogUseCase usecase.IInsertLogUseCase
	Channel    *amqp.Channel
}

type AMQPCTRLogHandler struct {
	LogUseCase usecase.IInsertCTRLogUseCase
	Channel    *amqp.Channel
}

type HttpLogHandler struct {
	ListUseCase usecase.IListLogsUseCase
}

func NewAMQPLogHandler(logUseCase usecase.IInsertLogUseCase, ch *amqp.Channel) *AMQPLogHandler {
	return &AMQPLogHandler{
		LogUseCase: logUseCase,
		Channel:    ch,
	}
}

func NewAMQPCTRLogHandler(logUseCase usecase.IInsertCTRLogUseCase, ch *amqp.Channel) *AMQPCTRLogHandler {
	return &AMQPCTRLogHandler{
		LogUseCase: logUseCase,
		Channel:    ch,
	}
}

func NewHttpLogHandler(listUseCase usecase.IListLogsUseCase) *HttpLogHandler {
	return &HttpLogHandler{
		ListUseCase: listUseCase,
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

// TODO: [Server] Improve the current RPC Implementation to reduct frontend delays
// https://github.com/okuda-seminar/log_service/issues/85
func (h *AMQPCTRLogHandler) HandleCTRLog(msg amqp.Delivery) {
	req, err := ParseAMQPCTRLog(msg)
	if err != nil {
		log.Println("failed to parse CTR log request:", err)
		// In case of invalid request, we should nack the message without requeue
		msg.Nack(false, false)
		return
	}
	logDto := &usecase.InsertCTRLogDto{
		EventType: req.EventType,
		CreatedAt: req.CreatedAt,
		ObjectID:  req.ObjectID,
	}
	err = h.LogUseCase.InsertCTRLog(context.Background(), logDto)
	if err != nil {
		log.Println("failed to insert CTR log:", err)
		// In case of internal server error, we should nack the message with requeue
		msg.Nack(false, true)
		return
	}

	// Acknowledge the message
	msg.Ack(false)
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

func ParseAMQPCTRLog(msg amqp.Delivery) (AMQPCTRLogRequest, error) {
	var req AMQPCTRLogRequest

	decoder := json.NewDecoder(bytes.NewReader(msg.Body))
	err := decoder.Decode(&req)
	if err != nil {
		return req, err
	}
	return req, err
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

func (h *HttpLogHandler) HandleLogList(w http.ResponseWriter, r *http.Request) {
	logs, err := h.ListUseCase.ListLogs(r.Context())
	if err != nil {
		log.Printf("Failed to list logs: %v", err)
		http.Error(w, fmt.Sprintf("Internal Server Error: %v", err), http.StatusInternalServerError)
		return
	}

	responseLogs := make([]HttpLogListResponse, len(logs))
	for i, eachLog := range logs {
		responseLogs[i] = HttpLogListResponse{
			LogLevel:           eachLog.LogLevel,
			Date:               eachLog.Date,
			DestinationService: eachLog.DestinationService,
			SourceService:      eachLog.SourceService,
			RequestType:        eachLog.RequestType,
			Content:            eachLog.Content,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(responseLogs); err != nil {
		log.Printf("Failed to encode logs: %v", err)
		http.Error(w, fmt.Sprintf("Internal Server Error: %v", err), http.StatusInternalServerError)
		return
	}

}
