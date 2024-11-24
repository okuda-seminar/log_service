package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	clientPresentation "log_service/internal/client/presentation"
	"log_service/internal/server/presentation"
)

type IInsertLogUseCase interface {
	Serve() error
}

type InsertLogUseCase struct {
	logPresentation clientPresentation.ILogPresentation
}

func NewInsertLogUseCase(logPresentation clientPresentation.ILogPresentation) *InsertLogUseCase {
	return &InsertLogUseCase{
		logPresentation: logPresentation,
	}
}

func (u *InsertLogUseCase) Serve(req presentation.AMQPLogRequest) error {
	msgs, qName, err := u.logPresentation.Consume()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	corrID := uuid.New().String()

	if err := u.logPresentation.Publish(ctx, qName, corrID, req); err != nil {
		return err
	}

	return u.logPresentation.Serve(msgs, corrID)
}
