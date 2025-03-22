package presentation

import (
	"bytes"
	"context"
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

func TestNewAMQPLogHandler(t *testing.T) {
	fixedTime := time.Date(2024, 9, 23, 23, 7, 32, 840757000, time.Local)
	logRequest, msg := testMsg(t, fixedTime)

	tests := []struct {
		name               string
		msg                amqp.Delivery
		mockFunc           func(m *usecase.MockIInsertLogUseCase)
		expectedStatusCode int
	}{
		{
			name: "success",
			msg:  msg,
			mockFunc: func(m *usecase.MockIInsertLogUseCase) {
				m.EXPECT().InsertLog(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, gotRequest *usecase.InsertLogDto) error {
					convertedRequest := AMQPLogRequest{
						LogLevel:           gotRequest.LogLevel,
						Date:               gotRequest.Date,
						DestinationService: gotRequest.DestinationService,
						SourceService:      gotRequest.SourceService,
						RequestType:        gotRequest.RequestType,
						Content:            gotRequest.Content,
					}
					testDiffLog(t, logRequest, convertedRequest)
					return nil
				}).Times(1)
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
				m.EXPECT().InsertLog(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, gotRequest *usecase.InsertLogDto) error {
					convertedRequest := AMQPLogRequest{
						LogLevel:           gotRequest.LogLevel,
						Date:               gotRequest.Date,
						DestinationService: gotRequest.DestinationService,
						SourceService:      gotRequest.SourceService,
						RequestType:        gotRequest.RequestType,
						Content:            gotRequest.Content,
					}
					testDiffLog(t, logRequest, convertedRequest)
					return errors.New("failed to insert log.")
				}).Times(1)
			},
			expectedStatusCode: utils.INTERNAL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

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
		testDiffLog(t, logRequest, req)
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
		testDiffCtrLog(t, ctrLogRequest, req)
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

func testDiffLog(t *testing.T, wantRequest AMQPLogRequest, gotRequest AMQPLogRequest) {

	if wantRequest.LogLevel != gotRequest.LogLevel {
		t.Errorf("Expected %s, got %s", wantRequest.LogLevel, gotRequest.LogLevel)
	}
	if !wantRequest.Date.Equal(gotRequest.Date) {
		t.Errorf("Expected %s, got %s", wantRequest.Date, gotRequest.Date)
	}
	if wantRequest.SourceService != gotRequest.SourceService {
		t.Errorf("Expected %s, got %s", wantRequest.SourceService, gotRequest.SourceService)
	}
	if wantRequest.DestinationService != gotRequest.DestinationService {
		t.Errorf("Expected %s, got %s", wantRequest.DestinationService, gotRequest.DestinationService)
	}
	if wantRequest.RequestType != gotRequest.RequestType {
		t.Errorf("Expected %s, got %s", wantRequest.RequestType, gotRequest.RequestType)
	}
	if wantRequest.Content != gotRequest.Content {
		t.Errorf("Expected %s, got %s", wantRequest.Content, gotRequest.Content)
	}
}

func testDiffCtrLog(t *testing.T, wantRequest AMQPCTRLogRequest, gotRequest AMQPCTRLogRequest) {

	if wantRequest.EventType != gotRequest.EventType {
		t.Errorf("Expected %s, got %s", wantRequest.EventType, gotRequest.EventType)
	}
	if !wantRequest.CreatedAt.Equal(gotRequest.CreatedAt) {
		t.Errorf("Expected %s, got %s", wantRequest.CreatedAt, gotRequest.CreatedAt)
	}
	if wantRequest.ObjectID != gotRequest.ObjectID {
		t.Errorf("Expected %s, got %s", wantRequest.ObjectID, gotRequest.ObjectID)
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
