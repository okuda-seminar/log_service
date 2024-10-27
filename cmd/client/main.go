package main

import (
	"flag"
	"log"
	"time"

	"log_service/internal/client"
	"log_service/internal/server/presentation"
)

func parseFlags() presentation.AMQPLogRequest {
	logLevel := flag.String("log-level", "INFO", "Log level")
	sourceService := flag.String("source-service", "", "Source service")
	destinationService := flag.String("destination-service", "", "Destination service")
	requestType := flag.String("request-type", "", "Request type")
	content := flag.String("content", "", "Log content")

	flag.Parse()

	req := presentation.AMQPLogRequest{
		LogLevel:           *logLevel,
		Date:               time.Now(),
		SourceService:      *sourceService,
		DestinationService: *destinationService,
		RequestType:        *requestType,
		Content:            *content,
	}

	return req
}

func main() {
	req := parseFlags()

	if err := client.Run(req); err != nil {
		log.Fatalf("Failed to send log request: %v", err)
	}
}
