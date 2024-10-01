package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"log_service/app/infrastructure/mysql/db"
	"log_service/app/infrastructure/mysql/repository"
	"log_service/app/infrastructure/rabbitmq"
	"log_service/app/presentation"
	"log_service/app/usecase"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	logRepo := repository.NewLogRepository(db)
	logUseCase := usecase.NewInsertLogUseCase(logRepo)
	amqpLogHandler := presentation.NewAMQPLogHandler(logUseCase)

	amqpConn, ch, msgs, err := rabbitmq.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	defer amqpConn.Close()
	defer ch.Close()

	done := make(chan bool)
	go func() {
		for d := range msgs {
			amqpLogHandler.HandleLog(d)
			log.Printf("Received a message: %s", d.Body)
			if err := d.Ack(false); err != nil {
				log.Fatalf("Failed to ack message: %v", err)
			}
		}
		done <- true
	}()

	log.Printf("Waiting for messages. To exit press CTRL^C")

	<-ctx.Done()
	stop()

	log.Println("received sigint/sigterm, shutting down...")
	log.Println("press Ctrl^C again to force shutdown")

	if err := ch.Cancel(rabbitmq.QUEUE_NAME, false); err != nil {
		log.Panic(err)
	}
	if err := ch.Close(); err != nil {
		log.Panic(err)
	}

	select {
	case <-done:
		log.Println("finished processing all jobs")
	case <-time.After(5 * time.Second):
		log.Println("timed out waiting for jobs to finish")
	}
}
