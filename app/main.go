package main

import (
	"log"
	"log_service/app/infrastructure/rabbitmq"
)

func main() {
	ch, msgs, err := rabbitmq.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer ch.Close()

	for msg := range msgs {
		log.Printf("Received message: %s", msg.Body)
	}
}
