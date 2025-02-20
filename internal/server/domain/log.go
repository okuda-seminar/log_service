package domain

import "time"

type Log struct {
	LogLevel           string
	Date               time.Time
	DestinationService string
	SourceService      string
	RequestType        string
	Content            string
}

// CTRLog represents a log entry for tracking user interactions with a page element.
//
// The zero value of CTRLog is not valid for use; values should be explicitly initialized.
type CTRLog struct {
	// EventType specifies the type of interaction, such as "click" or "impression".
	EventType string
	// CreatedAt is the timestamp when the event occurred.
	CreatedAt time.Time
	// ObjectID uniquely identifies the page element related to the event.
	ObjectID string
}

func NewLog(
	logLevel string,
	date time.Time,
	destinationService string,
	sourceService string,
	requestType string,
	content string,
) *Log {
	return &Log{
		LogLevel:           logLevel,
		Date:               date,
		DestinationService: destinationService,
		SourceService:      sourceService,
		RequestType:        requestType,
		Content:            content,
	}
}

// NewCTRLog creates a new CTRLog instance with the specified event type,
// creation timestamp, and associated object ID.
//
// eventType indicates the type of interaction, such as "click" or "impression".
// createdAt specifies the timestamp when the event occurred.
// objectid uniquely identifies the page element related to the event.
//
// Returns a pointer to a CTRLog instance initialized with the provided values.
func NewCTRLog(
	eventType string,
	createdAt time.Time,
	objectID string,
) *CTRLog {
	return &CTRLog{
		EventType: eventType,
		CreatedAt: createdAt,
		ObjectID:  objectID,
	}
}
