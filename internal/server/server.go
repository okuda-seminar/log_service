package server

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/rabbitmq/amqp091-go"

	"log_service/internal/server/infrastructure/di"
	"log_service/internal/server/infrastructure/rabbitmq"
	"log_service/internal/server/presentation"
)

func Run() error {
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	c, err := di.BuildLogContainer()
	if err != nil {
		return err
	}

	err = c.Invoke(func(
		dbConn *sql.DB,
		amqpConn *amqp091.Connection,
		amqpCh *amqp091.Channel,
		amqpMsgs <-chan amqp091.Delivery,
		amqpLogHandler *presentation.AMQPLogHandler,
		httpLogHander *presentation.HttpLogHandler,
	) {
		defer dbConn.Close()
		defer amqpCh.Close()
		defer amqpConn.Close()

		done := make(chan bool)
		go func() {
			for d := range amqpMsgs {
				amqpLogHandler.HandleLog(d)
				if err := d.Ack(false); err != nil {
					log.Fatalf("Failed to ack message: %v", err)
				}
			}
			done <- true
		}()

		mux := http.NewServeMux()
		mux.HandleFunc("/logs", httpLogHander.HandleLogList)

		srv := &http.Server{
			Addr:    ":8080",
			Handler: mux,
		}

		go func() {
			log.Printf("HTTP server is running on %s", srv.Addr)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Failed to start server: %v", err)
			}
		}()

		log.Printf("Waiting for messages. To exit press CTRL^C")

		<-ctx.Done()
		stop()

		log.Println("received sigint/sigterm, shutting down...")
		log.Println("press Ctrl^C again to force shutdown")

		if err := amqpCh.Cancel(rabbitmq.QUEUE_NAME, false); err != nil {
			log.Fatalf("Failed to cancel consumer: %v", err)
		}
		if err := amqpCh.Close(); err != nil {
			log.Fatalf("Failed to close channel: %v", err)
		}

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("HTTP server Shutdown: %v", err)
		}
		log.Println("HTTP server gracefully stopped")

		select {
		case <-done:
			log.Println("finished processing all jobs")
		case <-time.After(5 * time.Second):
			log.Println("timed out waiting for jobs to finish")
		}
	})

	if err != nil {
		return err
	}

	return nil
}
