include .env
export $(shell sed 's/=.*//' .env)

up:
	docker compose up -d

down:
	docker compose down -v --remove-orphans

migrate_up:
	migrate -path internal/server/infrastructure/mysql/db/schema -database "mysql://root:$(MYSQL_ROOT_PASSWORD)@tcp(localhost:3306)/${MYSQL_DATABASE}" up

migrate_down:
	migrate -path internal/server/infrastructure/mysql/db/schema -database "mysql://root:$(MYSQL_ROOT_PASSWORD)@tcp(localhost:3306)/${MYSQL_DATABASE}" down

exec_db:
	docker compose exec db mysql -u root -p$(MYSQL_ROOT_PASSWORD) ${MYSQL_DATABASE}

test:
	go test -cover ./... -coverprofile=cover.out
	go tool cover -html=cover.out

generate:
	sqlc generate

mock-gen:
	mockgen -package domain -source=internal/server/domain/log_repository.go -destination=internal/server/domain/log_mock.go
	mockgen -package usecase -source=internal/server/usecase/insert_log.go -destination=internal/server/usecase/insert_log_mock.go