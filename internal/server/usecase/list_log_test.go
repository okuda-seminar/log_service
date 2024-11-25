package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"log_service/internal/server/domain"
)

func TestListLog(t *testing.T) {
	t.Parallel()
	currTime := time.Now()

	testCases := map[string]struct {
		mockFunc  func(*domain.MockILogRepository)
		wantLogs  int
		wantError bool
	}{
		"ListLogs success": {
			mockFunc: func(m *domain.MockILogRepository) {
				sampleLog := &domain.Log{
					LogLevel:           "INFO",
					Date:               currTime,
					DestinationService: "UserService",
					SourceService:      "AuthService",
					RequestType:        "POST",
					Content:            "User created successfully.",
				}
				m.EXPECT().List(gomock.Any()).Return([]domain.Log{*sampleLog}, nil).Times(1)
			},
			wantLogs:  1,
			wantError: false,
		},
		"ListLogs failure": {
			mockFunc: func(m *domain.MockILogRepository) {
				m.EXPECT().List(gomock.Any()).Return(nil, errors.New("failed to list logs")).Times(1)
			},
			wantLogs:  0,
			wantError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			mockRepo := domain.NewMockILogRepository(ctrl)
			logListUseCase := NewListLogsUseCase(mockRepo)
			ctx := context.Background()
			tc.mockFunc(mockRepo)

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
