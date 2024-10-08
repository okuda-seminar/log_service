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
