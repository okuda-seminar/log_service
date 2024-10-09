package usecase

import (
	"context"
	"testing"
	"errors"
	"time"

	"go.uber.org/mock/gomock"

	"log_service/internal/server/domain"
)

func setupTest(t *testing.T) (*gomock.Controller, *domain.MockILogRepository, *ListLogsUseCase) {
    ctrl := gomock.NewController(t)
    mockRepo := domain.NewMockILogRepository(ctrl)
    logListUseCase := NewListLogsUseCase(mockRepo)
    return ctrl, mockRepo, logListUseCase
}


func TestListLog(t *testing.T) {
    _ , mockRepo, logListUseCase := setupTest(t)

    currTime := time.Now()

    testCases := map[string]struct {
        mockFunc   func()
        wantLogs   int
        wantError  bool
    }{
        "ListLogs success": {
            mockFunc: func() {
                sampleLog := &domain.Log{
                    LogLevel:           "INFO",
                    Date:               currTime,
                    DestinationService: "UserService",
                    SourceService:      "AuthService",
                    RequestType:        "POST",
                    Content:            "User created successfully.",
                }
                mockRepo.EXPECT().List(gomock.Any()).Return([]domain.Log{*sampleLog}, nil).Times(1)
            },
            wantLogs:  1,
            wantError: false,
        },
        "ListLogs failure": {
            mockFunc: func() {
                mockRepo.EXPECT().List(gomock.Any()).Return(nil, errors.New("failed to list logs")).Times(1)
            },
            wantLogs:  0,
            wantError: true,
        },
    }

    for name, tc := range testCases {
        t.Run(name, func(t *testing.T) {
			t.Parallel()
            ctx := context.Background()
            tc.mockFunc()

            results, err := logListUseCase.ListLogs(ctx)

            if tc.wantError {
                if err == nil {
                    t.Errorf("ListLogs() expected error but got none")
                }
            } else {
                if err != nil {
                    t.Errorf("ListLogs() unexpected error = %v", err)
                }
                if len(results) != tc.wantLogs {
                    t.Errorf("ListLogs() expected %d logs, got %d", tc.wantLogs, len(results))
                }
            }
        })
    }
}