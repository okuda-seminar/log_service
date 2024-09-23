package usecase

import (
	"context"
	"time"

	"log_service/app/domain"
)

type IInsertLogUseCase interface {
	InsertLog(ctx context.Context, dto *InsertLogDto) error
}

type InsertLogUseCase struct {
	logRepository domain.ILogRepository
}

func NewInsertLogUseCase(logRepository domain.ILogRepository) *InsertLogUseCase {
	return &InsertLogUseCase{
		logRepository: logRepository,
	}
}

type InsertLogDto struct {
	LogLevel           string
	Date               time.Time
	DestinationService string
	SourceService      string
	RequestType        string
	Content            string
}

func (u *InsertLogUseCase) InsertLog(ctx context.Context, dto *InsertLogDto) error {
	log := domain.NewLog(
		dto.LogLevel,
		dto.Date,
		dto.DestinationService,
		dto.SourceService,
		dto.RequestType,
		dto.Content,
	)
	return u.logRepository.Save(ctx, log)
}
