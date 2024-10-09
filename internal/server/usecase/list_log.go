package usecase

import (
	"context"
	"time"
	"log_service/internal/server/domain"
)

type IListLogsUseCase interface {
	ListLogs(ctx context.Context) ([]*domain.Log, error)
}

type ListLogsUseCase struct {
	logRepository domain.ILogRepository
}

func NewListLogsUseCase(logRepository domain.ILogRepository) *ListLogsUseCase {
	return &ListLogsUseCase{
		logRepository: logRepository,
	}
}

type ListLogsDto struct {
	LogLevel           string
	Date               time.Time
	DestinationService string
	SourceService      string
	RequestType        string
	Content            string
}

func (u *ListLogsUseCase) ListLogs(ctx context.Context) ([]*domain.Log, error) {
	logs, err := u.logRepository.List(ctx)
	if err != nil {
		return nil, err
	}
	var logPointers []*domain.Log
	for i := range logs {
		logPointers = append(logPointers, &logs[i])
	}

	return logPointers, nil
}

