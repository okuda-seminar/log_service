package presentation

import "time"

type AMQPLogRequest struct {
	LogLevel           string    `json:"log_level"`
	Date               time.Time `json:"date"`
	SourceService      string    `json:"source_service"`
	DestinationService string    `json:"destination_service"`
	RequestType        string    `json:"request_type"`
	Content            string    `json:"content"`
}

type AMQPCTRLogRequest struct {
	EventType string    `json:"eventType"`
	ObjectID  string    `json:"objectId"`
	CreatedAt time.Time `json:"createdAt"`
}
