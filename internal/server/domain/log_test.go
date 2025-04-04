package domain

import (
	"testing"
	"time"
)

func TestNewLog(t *testing.T) {
	// Setup test data
	logLevel := "INFO"
	date := time.Now()
	destinationService := "UserService"
	sourceService := "AuthService"
	requestType := "POST"
	content := "User created successfully."

	// Call the function
	log := NewLog(logLevel, date, destinationService, sourceService, requestType, content)

	// Check if the log is populated correctly
	if log.LogLevel != logLevel {
		t.Errorf("Expected LogLevel %s, got %s", logLevel, log.LogLevel)
	}
	if log.Date != date {
		t.Errorf("Expected Date %s, got %s", date, log.Date)
	}
	if log.DestinationService != destinationService {
		t.Errorf("Expected DestinationService %s, got %s", destinationService, log.DestinationService)
	}
	if log.SourceService != sourceService {
		t.Errorf("Expected SourceService %s, got %s", sourceService, log.SourceService)
	}
	if log.RequestType != requestType {
		t.Errorf("Expected RequestType %s, got %s", requestType, log.RequestType)
	}
	if log.Content != content {
		t.Errorf("Expected Content %s, got %s", content, log.Content)
	}
}

// TestNewCTRLog tests the NewCTRLog function
func TestNewCTRLog(t *testing.T) {
	// Setup test data
	eventType := "click"
	createdAt := time.Now()
	objectID := "123456"

	// Call the function
	ctrLog := NewCTRLog(eventType, createdAt, objectID)

	// Check if the log is populated correctly
	if ctrLog.CreatedAt != createdAt {
		t.Errorf("Expected CreatedAt %s, got %s", createdAt, ctrLog.CreatedAt)
	}
	if ctrLog.ObjectID != objectID {
		t.Errorf("Expected Objectid %s, got %s", objectID, ctrLog.ObjectID)
	}
	if ctrLog.EventType != eventType {
		t.Errorf("Expected EventType %s, got %s", eventType, ctrLog.EventType)
	}
}
