package repository

import (
	"context"
	"database/sql"

	"log_service/internal/server/domain"
	"log_service/internal/server/infrastructure/mysql/db/dbgen"
)

// LogRepository provides methods to interact with the log storage in the database.
type LogRepository struct {
	db *sql.DB
}

// NewLogRepository creates a new instance of LogRepository with the given database connection.
func NewLogRepository(db *sql.DB) *LogRepository {
	return &LogRepository{
		db: db,
	}
}

// Save stores a new log entry into the database.
// It takes a context and a Log object from the domain package as arguments.
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

// List retrieves all log entries from the database.
// It returns a slice of Log objects from the domain package or an error if the query fails.
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

// CTRSave stores a new CTRLog entry into the database.
// It takes a context and a CTRLog object from the domain package as arguments.
func (r *LogRepository) CTRSave(ctx context.Context, ctrLog *domain.CTRLog) error {
	err := dbgen.New(r.db).InsertCTRLog(ctx, dbgen.InsertCTRLogParams{
		EventType: ctrLog.EventType,
		CreatedAt: ctrLog.CreatedAt,
		ObjectID:  ctrLog.ObjectID,
	})
	return err
}

// CTRList retrieves all CTRLog entries from the database.
// It returns a slice of CTRLog objects from the domain package or an error if the query fails.
func (r *LogRepository) CTRList(ctx context.Context) ([]domain.CTRLog, error) {
	ctrLogs, err := dbgen.New(r.db).ListCTRLogs(ctx)
	if err != nil {
		return nil, err
	}

	var result []domain.CTRLog
	for _, ctrLog := range ctrLogs {
		result = append(result, domain.CTRLog{
			EventType: ctrLog.EventType,
			CreatedAt: ctrLog.CreatedAt,
			ObjectID:  ctrLog.ObjectID,
		})
	}

	return result, nil
}
