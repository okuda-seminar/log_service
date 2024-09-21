package repository

import (
	"context"
	"log_service/app/domain"
	"testing"
	"time"
)

func TestInsertLog(t *testing.T) {
	repo := NewLogRepository(dbConnTest)
	err := repo.Save(context.Background(), &domain.Log{
		LogLevel:           "INFO",
		Date:               time.Now(),
		DestinationService: "UserService",
		SourceService:      "AuthService",
		RequestType:        "POST",
		Content:            "User created successfully.",
	})
	if err != nil {
		t.Error(err)
	}
}
