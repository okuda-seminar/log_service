// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: log.sql

package dbgen

import (
	"context"
	"time"
)

const insertCTRLog = `-- name: InsertCTRLog :exec
INSERT INTO ctr_logs (
  event_type, created_at, object_id
) VALUES (
  ?, ?, ?
)
`

type InsertCTRLogParams struct {
	EventType string
	CreatedAt time.Time
	ObjectID  string
}

func (q *Queries) InsertCTRLog(ctx context.Context, arg InsertCTRLogParams) error {
	_, err := q.db.ExecContext(ctx, insertCTRLog, arg.EventType, arg.CreatedAt, arg.ObjectID)
	return err
}

const insertLog = `-- name: InsertLog :exec
INSERT INTO logs (
  log_level, date, destination_service, source_service, request_type, content
) VALUES (
  ?, ?, ?, ?, ?, ?
)
`

type InsertLogParams struct {
	LogLevel           string
	Date               time.Time
	DestinationService string
	SourceService      string
	RequestType        string
	Content            string
}

func (q *Queries) InsertLog(ctx context.Context, arg InsertLogParams) error {
	_, err := q.db.ExecContext(ctx, insertLog,
		arg.LogLevel,
		arg.Date,
		arg.DestinationService,
		arg.SourceService,
		arg.RequestType,
		arg.Content,
	)
	return err
}

const listCTRLogs = `-- name: ListCTRLogs :many
SELECT
  event_type, created_at, object_id
FROM ctr_logs
`

func (q *Queries) ListCTRLogs(ctx context.Context) ([]CtrLog, error) {
	rows, err := q.db.QueryContext(ctx, listCTRLogs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CtrLog
	for rows.Next() {
		var i CtrLog
		if err := rows.Scan(&i.EventType, &i.CreatedAt, &i.ObjectID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listLogs = `-- name: ListLogs :many
SELECT
  log_level, date, destination_service, source_service, request_type, content
FROM logs
`

func (q *Queries) ListLogs(ctx context.Context) ([]Log, error) {
	rows, err := q.db.QueryContext(ctx, listLogs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Log
	for rows.Next() {
		var i Log
		if err := rows.Scan(
			&i.LogLevel,
			&i.Date,
			&i.DestinationService,
			&i.SourceService,
			&i.RequestType,
			&i.Content,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
