-- name: InsertLog :exec
INSERT INTO logs (
  log_level, date, destination_service, source_service, request_type, content
) VALUES (
  ?, ?, ?, ?, ?, ?
);

-- name: ListLogs :many
SELECT
  log_level, date, destination_service, source_service, request_type, content
FROM logs
;