package main

import (
	"log"

	"log_service/app/infrastructure/mysql/db"
	"log_service/app/infrastructure/mysql/repository"
	"log_service/app/infrastructure/rabbitmq"
	"log_service/app/presentation"
	"log_service/app/usecase"
)

func main() {
	ch, msgs, err := rabbitmq.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer ch.Close()

	db, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	logRepo := repository.NewLogRepository(db)
	logUseCase := usecase.NewInsertLogUseCase(logRepo)
	amqpLogHandler := presentation.NewAMQPLogHandler(logUseCase)

	for msg := range msgs {
		amqpLogHandler.HandleLog(msg)
	}
}
