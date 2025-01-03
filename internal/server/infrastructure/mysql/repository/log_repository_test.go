package repository

import (
	"context"
	"testing"
	"time"

	"log_service/internal/server/domain"
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

// TODO: https://github.com/okuda-seminar/log_service/issues/18#issue-2557594866
// - [Server] Improve Test Cleanup Phase to Ensure Proper Test Isolation
func TestList(t *testing.T) {
	repo := NewLogRepository(dbConnTest)
	repo.Save(context.Background(), &domain.Log{
		LogLevel:           "INFO",
		Date:               time.Now(),
		DestinationService: "UserService",
		SourceService:      "AuthService",
		RequestType:        "POST",
		Content:            "Test Get Log.",
	})
	results, err := repo.List(context.Background())
	if err != nil {
		t.Error(err)
	}
	if len(results) != 2 {
		t.Errorf("Expected 2 log, got %d", len(results))
	}
}

func TestInsertCTRLog(t *testing.T) {
	repo := NewLogRepository(dbConnTest)
	err := repo.CTRSave(context.Background(), &domain.CTRLog{
		CreatedAt: time.Now(),
		Objectid:  "123456",
	})
	if err != nil {
		t.Error("Failed to insert CTR log:", err)
	}
}
