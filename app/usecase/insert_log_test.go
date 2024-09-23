package usecase

import (
	"context"
	"log_service/app/domain"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
)

func TestInsertLog(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUserRepo := domain.NewMockILogRepository(ctrl)
	logInsertUseCase := NewInsertLogUseCase(mockUserRepo)

	currTime := time.Now()
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
				mockUserRepo.EXPECT().Save(
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
			ctx := context.Background()
			tt.mockFunc()
			err := logInsertUseCase.InsertLog(ctx, tt.dto)
			if err != nil {
				t.Errorf("InsertLog() error = %v", err)
			}
		})
	}
}
