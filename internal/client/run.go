package client

import (
	"github.com/google/uuid"

	"log_service/internal/client/infrastructure/rabbitmq"
	clientPresentation "log_service/internal/client/presentation"
	"log_service/internal/client/usecase"
	"log_service/internal/server/presentation"
)

func Run(req presentation.AMQPLogRequest) error {
	conn, ch, err := rabbitmq.Connect()
	if err != nil {
		return err
	}
	defer conn.Close()
	defer ch.Close()

	corrID := uuid.New().String()

	logPresentation := clientPresentation.NewLogPresentation(ch, corrID)
	logUseCase := usecase.NewInsertLogUseCase(logPresentation)

	return logUseCase.Serve(req)
}
