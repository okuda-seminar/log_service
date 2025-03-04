package usecase

import (
	"context"
	"time"

	"log_service/internal/server/domain"
)

type IInsertLogUseCase interface {
	InsertLog(ctx context.Context, dto *InsertLogDto) error
}

// IInsertCTRLogUseCase is an interface for inserting CTR logs.
type IInsertCTRLogUseCase interface {
	InsertCTRLog(ctx context.Context, dto *InsertCTRLogDto) error
}

type InsertLogUseCase struct {
	logRepository domain.ILogRepository
}

// InsertLogUseCase is a use case for inserting log entries into the database.
type InsertCTRLogUseCase struct {
	logRepository domain.ILogRepository
}

func NewInsertLogUseCase(logRepository domain.ILogRepository) *InsertLogUseCase {
	return &InsertLogUseCase{
		logRepository: logRepository,
	}
}

// NewInsertCTRLogUseCase creates a new instance of InsertCTRLogUseCase with the given log repository.
func NewInsertCTRLogUseCase(logRepository domain.ILogRepository) *InsertCTRLogUseCase {
	return &InsertCTRLogUseCase{
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

// InsertCTRLogDto is a data transfer object for inserting CTR logs.
type InsertCTRLogDto struct {
	EventType string
	CreatedAt time.Time
	ObjectID  string
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

// InsertCTRLog inserts a new CTR log entry into the database.
// It takes a context and a CTRLogDto object as arguments.
// It returns an error if the operation fails.
func (u *InsertLogUseCase) InsertCTRLog(ctx context.Context, dto *InsertCTRLogDto) error {
	ctrLog := domain.NewCTRLog(
		dto.EventType,
		dto.CreatedAt,
		dto.ObjectID,
	)
	return u.logRepository.CTRSave(ctx, ctrLog)
}
