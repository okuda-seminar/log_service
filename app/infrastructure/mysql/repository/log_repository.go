package repository

import (
	"context"
	"database/sql"
	"log_service/app/domain"
	"log_service/app/infrastructure/mysql/db/dbgen"
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
