MYSQL_ROOT_PASSWORD=$(shell cat db/password.txt)

up:
	docker compose up -d

down:
	docker compose down -v --remove-orphans

migrate_up:
	migrate -path app/infrastructure/mysql/db/schema -database "mysql://root:password@tcp(localhost:3306)/example" up

migrate_down:
	migrate -path app/infrastructure/mysql/db/schema -database "mysql://root:password@tcp(localhost:3306)/example" down

exec_db:
	docker compose exec db mysql -u root -p$(MYSQL_ROOT_PASSWORD) example

test:
	go test -cover ./... -coverprofile=cover.out
	go tool cover -html=cover.out

generate:
	sqlc generate

mock-gen:
	mockgen -package domain -source=app/domain/log_repository.go -destination=app/domain/log_mock.go
	mockgen -package usecase -source=app/usecase/insert_log.go -destination=app/usecase/insert_log_mock.go