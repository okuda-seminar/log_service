package main

import (
	"encoding/json"
	"flag"
	"log"
	"time"

	"log_service/internal/client/infrastructure/rabbitmq"
	"log_service/internal/server/presentation"
)

func parseFlags() *presentation.AMQPLogRequest {
	logLevel := flag.String("log-level", "INFO", "Log level")
	sourceService := flag.String("source-service", "", "Source service")
	destinationService := flag.String("destination-service", "", "Destination service")
	requestType := flag.String("request-type", "", "Request type")
	content := flag.String("content", "", "Log content")

	flag.Parse()

	req := &presentation.AMQPLogRequest{
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
	conn, ch, err := rabbitmq.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	defer conn.Close()
	defer ch.Close()

	req := parseFlags()

	var bytes []byte
	bytes, err = json.Marshal(req)
	if err != nil {
		log.Fatalf("Failed to marshal log request: %v", err)
	}
	if err := rabbitmq.Publish(ch, bytes); err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}

	log.Printf("Sent a message: %s", string(bytes))
}
