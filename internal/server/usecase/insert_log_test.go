package usecase

import (
	"context"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"log_service/internal/server/domain"
)

func TestInsertLog(t *testing.T) {
	t.Parallel()

	currTime := time.Now()
	tests := []struct {
		name     string
		dto      *InsertLogDto
		mockFunc func(*domain.MockILogRepository)
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
			mockFunc: func(m *domain.MockILogRepository) {
				m.EXPECT().Save(
					gomock.Any(),
					&domain.Log{
						LogLevel:           "INFO",
						Date:               currTime,
						DestinationService: "UserService",
						SourceService:      "AuthService",
						RequestType:        "POST",
						Content:            "User created successfully.",
					},
				).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			mockUserRepo := domain.NewMockILogRepository(ctrl)
			logInsertUseCase := NewInsertLogUseCase(mockUserRepo)
			ctx := context.Background()
			tt.mockFunc(mockUserRepo)
			err := logInsertUseCase.InsertLog(ctx, tt.dto)
			if err != nil {
				t.Errorf("InsertLog() error = %v", err)
			}
		})
	}
}
