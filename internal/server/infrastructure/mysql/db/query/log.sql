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

-- name: InsertCTRLog :exec
INSERT INTO ctr_logs (
  event_type, created_at, object_id
) VALUES (
  ?, ?, ?
);

-- name: ListCTRLogs :many
SELECT
  event_type, created_at, object_id
FROM ctr_logs
;