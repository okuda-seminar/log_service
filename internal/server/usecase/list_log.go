package usecase

import (
	"context"
	"time"

	"log_service/internal/server/domain"
)

type IListLogsUseCase interface {
	ListLogs(ctx context.Context) ([]*ListLogDto, error)
}

type ListLogsUseCase struct {
	logRepository domain.ILogRepository
}

func NewListLogsUseCase(logRepository domain.ILogRepository) *ListLogsUseCase {
	return &ListLogsUseCase{
		logRepository: logRepository,
	}
}

// TODO: [Server] Implement LogID Assignment for Logs
// https://github.com/okuda-seminar/log_service/issues/34#issue-2601726013
type ListLogDto struct {
	LogLevel           string
	Date               time.Time
	DestinationService string
	SourceService      string
	RequestType        string
	Content            string
}

func (u *ListLogsUseCase) ListLogs(ctx context.Context) ([]*ListLogDto, error) {
	logs, err := u.logRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	var logDtos []*ListLogDto
	for _, log := range logs {
		logDto := &ListLogDto{
			LogLevel:           log.LogLevel,
			Date:               log.Date,
			DestinationService: log.DestinationService,
			SourceService:      log.SourceService,
			RequestType:        log.RequestType,
			Content:            log.Content,
		}
		logDtos = append(logDtos, logDto)
	}
	return logDtos, nil
}
