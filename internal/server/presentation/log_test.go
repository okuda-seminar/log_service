package presentation

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/google/go-cmp/cmp"
	amqp "github.com/rabbitmq/amqp091-go"

	"go.uber.org/mock/gomock"

	"log_service/internal/server/usecase"
	"log_service/internal/utils"
)

type errorWriterResponse struct {
	statusCode int
	body       string
}

func (e *errorWriterResponse) Header() http.Header {
	return http.Header{}
}

func (e *errorWriterResponse) Write(p []byte) (int, error) {
	e.body = string(p)
	return 0, errors.New("failed to encode logs")
}

func (e *errorWriterResponse) WriteHeader(statusCode int) {
	e.statusCode = statusCode
}

func TestNewAMQPLogHandler(t *testing.T) {
	fixedTime := time.Date(2024, 9, 23, 23, 7, 32, 840757000, time.Local)
	logRequest, msg := testMsg(t, fixedTime)
	tests := []struct {
		name               string
		msg                amqp.Delivery
		mockFunc           func(m *usecase.MockIInsertLogUseCase)
		expectedStatusCode int
		expectedMessage    string
	}{
		{
			name: "success",
			msg:  msg,
			mockFunc: func(m *usecase.MockIInsertLogUseCase) {
				m.EXPECT().InsertLog(
					gomock.Any(),
					&usecase.InsertLogDto{
						LogLevel:           logRequest.LogLevel,
						Date:               logRequest.Date,
						DestinationService: logRequest.DestinationService,
						SourceService:      logRequest.SourceService,
						RequestType:        logRequest.RequestType,
						Content:            logRequest.Content,
					},
				).Times(1).Return(nil)
			},
			expectedStatusCode: utils.OK,
		},
		{
			name:               "failed",
			msg:                amqp.Delivery{},
			mockFunc:           func(m *usecase.MockIInsertLogUseCase) {},
			expectedStatusCode: utils.INVALID_ARGUMENT,
		},
		{
			name: "failed",
			msg:  msg,
			mockFunc: func(m *usecase.MockIInsertLogUseCase) {
				m.EXPECT().InsertLog(
					gomock.Any(),
					&usecase.InsertLogDto{
						LogLevel:           logRequest.LogLevel,
						Date:               logRequest.Date,
						DestinationService: logRequest.DestinationService,
						SourceService:      logRequest.SourceService,
						RequestType:        logRequest.RequestType,
						Content:            logRequest.Content,
					},
				).Times(1).Return(errors.New("failed to insert log."))
			},
			expectedStatusCode: utils.INTERNAL,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockInsertUseCase := usecase.NewMockIInsertLogUseCase(ctrl)
			handler := NewAMQPLogHandler(mockInsertUseCase, &amqp.Channel{})
			tt.mockFunc(mockInsertUseCase)
			var patchResponseCode int
			// gomonkey cannot be used for parallel tests because it operates on shared resources.
			patch := gomonkey.ApplyMethod(
				reflect.TypeOf(&amqp.Channel{}),
				"Publish",
				func(
					_ *amqp.Channel,
					_ string,
					_ string,
					_ bool,
					_ bool,
					msg amqp.Publishing,
				) error {
					res := &AmqpLogResponse{}
					err := json.Unmarshal(msg.Body, res)
					if err != nil {
						return err
					}
					patchResponseCode = res.StatusCode
					return nil
				})
			handler.HandleLog(tt.msg)
			t.Log(patchResponseCode)

			if patchResponseCode != tt.expectedStatusCode {
				t.Errorf("Expected %d, got %d", tt.expectedStatusCode, patchResponseCode)
			}
			patch.Reset()
		})
	}
}

type MockDelivery struct {
	amqp.Delivery
	ctrl     *gomock.Controller
	ackFunc  func(multiple bool) error
	nackFunc func(multiple, requeue bool) error
}

func NewMockDelivery(ctrl *gomock.Controller, body []byte) *MockDelivery {
	return &MockDelivery{
		Delivery: amqp.Delivery{Body: body},
		ctrl:     ctrl,
	}
}

func (m *MockDelivery) Ack(multiple bool) error {
	if m.ackFunc != nil {
		return m.ackFunc(multiple)
	}
	return nil
}

func (m *MockDelivery) Nack(multiple, requeue bool) error {
	if m.nackFunc != nil {
		return m.nackFunc(multiple, requeue)
	}
	return nil
}

// func TestNewAMQPLogHandler(t *testing.T) {
func TestNewAMQPCTRLogHandler(t *testing.T) {
	fixedTime := time.Date(2024, 9, 23, 23, 7, 32, 840757000, time.UTC)
	ctrLogRequest, msg := testCTRMsg(t, fixedTime)

	tests := []struct {
		name          string
		msg           *MockDelivery
		mockFunc      func(m *usecase.MockIInsertCTRLogUseCase)
		expectAck     bool
		expectNack    bool
		requeueOnNack bool
	}{
		{
			name: "success",
			msg:  NewMockDelivery(nil, msg.Body),
			mockFunc: func(m *usecase.MockIInsertCTRLogUseCase) {
				m.EXPECT().InsertCTRLog(
					gomock.Any(),
					&usecase.InsertCTRLogDto{
						EventType: ctrLogRequest.EventType,
						CreatedAt: ctrLogRequest.CreatedAt,
						ObjectID:  ctrLogRequest.ObjectID,
					},
				).Return(nil)
			},
			expectAck:  true,
			expectNack: false,
		},
		{
			name:          "failed parse",
			msg:           NewMockDelivery(nil, []byte("")),
			mockFunc:      func(m *usecase.MockIInsertCTRLogUseCase) {},
			expectAck:     false,
			expectNack:    true,
			requeueOnNack: false,
		},
		{
			name: "failed insert",
			msg:  NewMockDelivery(nil, msg.Body),
			mockFunc: func(m *usecase.MockIInsertCTRLogUseCase) {
				m.EXPECT().InsertCTRLog(
					gomock.Any(),
					&usecase.InsertCTRLogDto{
						EventType: "event",
						CreatedAt: fixedTime,
						ObjectID:  "1234",
					},
				).Return(errors.New("failed to insert CTR log."))
			},
			expectAck:     false,
			expectNack:    true,
			requeueOnNack: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUseCase := usecase.NewMockIInsertCTRLogUseCase(ctrl)
			tt.mockFunc(mockUseCase)

			mockMsg := NewMockDelivery(ctrl, tt.msg.Body)
			mockMsg.ackFunc = func(multiple bool) error {
				if tt.expectAck {
					return nil
				}
				t.Errorf("Ack should not have been called")
				return errors.New("unexpected ack call")
			}
			mockMsg.nackFunc = func(multiple, requeue bool) error {
				if tt.expectNack {
					if requeue != tt.requeueOnNack {
						t.Errorf("Expected Nack requeue=%v, got %v", tt.requeueOnNack, requeue)
						return errors.New("unexpected nack requeue value")
					}
					return nil
				}
				t.Errorf("Nack should not have been called")
				return errors.New("unexpected nack call")
			}

			handler := NewAMQPCTRLogHandler(mockUseCase, &amqp.Channel{})
			handler.HandleCTRLog(mockMsg.Delivery)

			if tt.expectAck {
				if mockMsg.ackFunc == nil {
					t.Errorf("Expected Ack() to be called, but it was not")
				}
			}
			if tt.expectNack {
				if mockMsg.nackFunc == nil {
					t.Errorf("Expected Nack() to be called, but it was not")
				}
			}
		})
	}
}

func TestParseAMQPLog(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
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
	})

	t.Run("failed", func(t *testing.T) {
		t.Parallel()
		msg := amqp.Delivery{}
		_, err := ParseAMQPLog(msg)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestParseAMQPCTRLog(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		fixedTime := time.Date(2024, 9, 23, 23, 7, 32, 840757000, time.Local)
		ctrLogRequest, msg := testCTRMsg(t, fixedTime)

		req, err := ParseAMQPCTRLog(msg)
		if err != nil {
			t.Error(err)
		}
		if req.EventType != ctrLogRequest.EventType {
			t.Errorf("Expected %s, got %s", ctrLogRequest.EventType, req.EventType)
		}
		if req.CreatedAt != ctrLogRequest.CreatedAt {
			t.Errorf("Expected %s, got %s", ctrLogRequest.CreatedAt, req.CreatedAt)
		}
		if req.ObjectID != ctrLogRequest.ObjectID {
			t.Errorf("Expected %s, got %s", ctrLogRequest.ObjectID, req.ObjectID)
		}
	})
	t.Run("failed", func(t *testing.T) {
		t.Parallel()
		msg := amqp.Delivery{}
		_, err := ParseAMQPCTRLog(msg)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
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

func testCTRMsg(t *testing.T, fixedTime time.Time) (AMQPCTRLogRequest, amqp.Delivery) {
	logRequest := AMQPCTRLogRequest{
		EventType: "event",
		CreatedAt: fixedTime,
		ObjectID:  "1234",
	}

	var payload bytes.Buffer
	err := json.NewEncoder(&payload).Encode(logRequest)
	if err != nil {
		t.Error(err)
	}

	msg := amqp.Delivery{Body: payload.Bytes()}
	return logRequest, msg
}

func SetupLogListTest(t *testing.T) (*gomock.Controller, *usecase.MockIListLogsUseCase, *HttpLogHandler) {
	ctrl := gomock.NewController(t)
	mockListUseCase := usecase.NewMockIListLogsUseCase(ctrl)
	handler := NewHttpLogHandler(mockListUseCase)
	return ctrl, mockListUseCase, handler
}

func TestHandleLogList(t *testing.T) {
	t.Parallel()
	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		_, mockListUseCase, handler := SetupLogListTest(t)

		now := time.Now()
		expectedLogs := []*usecase.ListLogDto{
			{
				LogLevel:           "INFO",
				Date:               now,
				DestinationService: "ServiceA",
				SourceService:      "ServiceB",
				RequestType:        "GET",
				Content:            "First log message",
			},
			{
				LogLevel:           "ERROR",
				Date:               now,
				DestinationService: "ServiceC",
				SourceService:      "ServiceD",
				RequestType:        "POST",
				Content:            "Second log message",
			},
		}

		mockListUseCase.EXPECT().ListLogs(gomock.Any()).Return(expectedLogs, nil).Times(1)

		req, err := http.NewRequest("GET", "/logs", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		handler.HandleLogList(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expectedResponse := []HttpLogListResponse{
			{
				LogLevel:           "INFO",
				Date:               now,
				DestinationService: "ServiceA",
				SourceService:      "ServiceB",
				RequestType:        "GET",
				Content:            "First log message",
			},
			{
				LogLevel:           "ERROR",
				Date:               now,
				DestinationService: "ServiceC",
				SourceService:      "ServiceD",
				RequestType:        "POST",
				Content:            "Second log message",
			},
		}

		var got []HttpLogListResponse
		if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
			t.Fatalf("Failed to decode response body: %v", err)
		}

		if diff := cmp.Diff(expectedResponse, got); diff != "" {
			t.Errorf("handler returned unexpected JSON (-want +got):\n%s", diff)
		}
	})

	t.Run("ListLogs Failure", func(t *testing.T) {
		t.Parallel()
		_, mockListUseCase, handler := SetupLogListTest(t)

		mockListUseCase.EXPECT().ListLogs(gomock.Any()).Return(nil, errors.New("failed to list logs")).Times(1)

		req, err := http.NewRequest("GET", "/logs", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		handler.HandleLogList(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}

		expectedError := "Internal Server Error: failed to list logs\n"
		if rr.Body.String() != expectedError {
			t.Errorf("handler returned unexpected error message: got %v want %v", rr.Body.String(), expectedError)
		}
	})

	t.Run("Encode Failure", func(t *testing.T) {
		t.Parallel()
		_, mockListUseCase, handler := SetupLogListTest(t)

		time := time.Now()
		logs := []*usecase.ListLogDto{
			{
				LogLevel:           "INFO",
				Date:               time,
				DestinationService: "ServiceA",
				SourceService:      "ServiceB",
				RequestType:        "GET",
				Content:            "First log message",
			},
		}
		mockListUseCase.EXPECT().ListLogs(gomock.Any()).Return(logs, nil).Times(1)

		errorWriter := &errorWriterResponse{}

		req, err := http.NewRequest("GET", "/logs", nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.HandleLogList(errorWriter, req)

		if errorWriter.statusCode != http.StatusInternalServerError {
			t.Errorf(
				"handler returned wrong status code: got %v want %v",
				errorWriter.statusCode,
				http.StatusInternalServerError,
			)
		}

		expected := "Internal Server Error: failed to encode logs\n"
		if errorWriter.body != expected {
			t.Errorf("handler returned unexpected error message: got %v want %v", errorWriter.body, expected)
		}
	})
}
