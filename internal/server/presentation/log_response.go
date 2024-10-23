package presentation

import "time"

type AmqpLogResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type HttpLogListResponse struct {
	LogLevel           string    `json:"log_level"`
	Date               time.Time `json:"date"`
	SourceService      string    `json:"source_service"`
	DestinationService string    `json:"destination_service"`
	RequestType        string    `json:"request_type"`
	Content            string    `json:"content"`
}
