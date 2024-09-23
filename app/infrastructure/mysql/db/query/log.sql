-- name: InsertLog :exec
INSERT INTO logs (
  log_level, date, destination_service, source_service, request_type, content
) VALUES (
  ?, ?, ?, ?, ?, ?
);