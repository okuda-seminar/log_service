package repository

import (
	"context"
	"database/sql"

	"log_service/internal/server/domain"
	"log_service/internal/server/infrastructure/mysql/db/dbgen"
)

type LogRepository struct {
	db *sql.DB
}

func NewLogRepository(db *sql.DB) *LogRepository {
	return &LogRepository{
		db: db,
	}
}

func (r *LogRepository) Save(ctx context.Context, log *domain.Log) error {
	err := dbgen.New(r.db).InsertLog(ctx, dbgen.InsertLogParams{
		LogLevel:           log.LogLevel,
		Date:               log.Date,
		DestinationService: log.DestinationService,
		SourceService:      log.SourceService,
		RequestType:        log.RequestType,
		Content:            log.Content,
	})
	return err
}

func (r *LogRepository) List(ctx context.Context) ([]domain.Log, error) {
	logs, err := dbgen.New(r.db).ListLogs(ctx)
	if err != nil {
		return nil, err
	}

	var result []domain.Log
	for _, log := range logs {
		result = append(result, domain.Log{
			LogLevel:           log.LogLevel,
			Date:               log.Date,
			DestinationService: log.DestinationService,
			SourceService:      log.SourceService,
			RequestType:        log.RequestType,
			Content:            log.Content,
		})
	}

	return result, nil
}
