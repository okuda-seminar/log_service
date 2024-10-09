package usecase

import (
	"context"
	"testing"
	"time"
	"go.uber.org/mock/gomock"
	"log_service/internal/server/domain"
)

func TestListLog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := domain.NewMockILogRepository(ctrl)
	logListUseCase := NewListLogsUseCase(mockRepo)
	logInsertUseCase := NewInsertLogUseCase(mockRepo)

	currTime := time.Now()

	// ログのサンプルデータ
	sampleLog := &domain.Log{
		LogLevel:           "INFO",
		Date:               currTime,
		DestinationService: "UserService",
		SourceService:      "AuthService",
		RequestType:        "POST",
		Content:            "User created successfully.",
	}

	tests := []struct {
		name     string
		dto      *InsertLogDto
		mockFunc func()
	}{
		{
			name: "success",
			dto: &InsertLogDto{
				LogLevel:           "INFO",
				Date:               currTime,
				DestinationService: "UserService",
				SourceService:      "AuthService",
				RequestType:        "POST",
				Content:            "User created successfully.",
			},
			mockFunc: func() {
				mockRepo.EXPECT().Save(gomock.Any(), sampleLog).Return(nil).Times(1)
			},
		},
		{
			name: "list success",
			mockFunc: func() {
				mockRepo.EXPECT().List(gomock.Any()).Return([]domain.Log{*sampleLog}, nil).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.dto != nil {
				tt.mockFunc()
				err := logInsertUseCase.InsertLog(ctx, tt.dto)
				if err != nil {
					t.Errorf("InsertLog() error = %v", err)
				}
			} else {
				tt.mockFunc()
			}
		})
	}

	// Test ListLogs use case
	results, err := logListUseCase.ListLogs(context.Background())
	if err != nil {
		t.Error(err)
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 log, got %d", len(results))
	}
}
