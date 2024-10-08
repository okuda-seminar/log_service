package presentation

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"go.uber.org/mock/gomock"

	"log_service/internal/server/usecase"
)

func TestNewAMQPLogHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInsertUseCase := usecase.NewMockIInsertLogUseCase(ctrl)
	handler := NewAMQPLogHandler(mockInsertUseCase)

	fixedTime := time.Date(2024, 9, 23, 23, 7, 32, 840757000, time.Local)
	logRequest, msg := testMsg(t, fixedTime)
	tests := []struct {
		name     string
		msg      amqp.Delivery
		mockFunc func()
	}{
		{
			name: "success",
			msg:  msg,
			mockFunc: func() {
				mockInsertUseCase.EXPECT().InsertLog(
					gomock.Any(),
					&usecase.InsertLogDto{
						LogLevel:           logRequest.LogLevel,
						Date:               logRequest.Date,
						DestinationService: logRequest.DestinationService,
						SourceService:      logRequest.SourceService,
						RequestType:        logRequest.RequestType,
						Content:            logRequest.Content,
					},
				).Times(1)
			},
		},
		{
			name:     "failed",
			msg:      amqp.Delivery{},
			mockFunc: func() {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			handler.HandleLog(tt.msg)
		})
	}
}

func TestParseAMQPLog(t *testing.T) {
	fixedTime := time.Date(2024, 9, 23, 23, 7, 32, 840757000, time.Local)
	logRequest, msg := testMsg(t, fixedTime)

	req, err := ParseAMQPLog(msg)
	if err != nil {
		t.Error(err)
	}
	if req.LogLevel != logRequest.LogLevel {
		t.Errorf("Expected %s, got %s", logRequest.LogLevel, req.LogLevel)
	}
	if req.Date != logRequest.Date {
		t.Errorf("Expected %s, got %s", logRequest.Date, req.Date)
	}
	if req.SourceService != logRequest.SourceService {
		t.Errorf("Expected %s, got %s", logRequest.SourceService, req.SourceService)
	}
	if req.DestinationService != logRequest.DestinationService {
		t.Errorf("Expected %s, got %s", logRequest.DestinationService, req.DestinationService)
	}
	if req.RequestType != logRequest.RequestType {
		t.Errorf("Expected %s, got %s", logRequest.RequestType, req.RequestType)
	}
	if req.Content != logRequest.Content {
		t.Errorf("Expected %s, got %s", logRequest.Content, req.Content)
	}
}

func testMsg(t *testing.T, fixedTime time.Time) (AMQPLogRequest, amqp.Delivery) {
	logRequest := AMQPLogRequest{
		LogLevel:           "INFO",
		Date:               fixedTime,
		SourceService:      "UserService",
		DestinationService: "AuthService",
		RequestType:        "POST",
		Content:            "User created successfully.",
	}

	var payload bytes.Buffer
	err := json.NewEncoder(&payload).Encode(logRequest)
	if err != nil {
		t.Error(err)
	}

	msg := amqp.Delivery{Body: payload.Bytes()}
	return logRequest, msg
}
