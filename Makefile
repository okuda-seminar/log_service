include .env
export $(shell sed 's/=.*//' .env)
GENERATE_IMAGE ?= generate-mock

ARGS ?= --log-level=debug --source-service=client --destination-service=server --request-type=POST --content="Hello World!"

up:
	docker compose up -d --build

down:
	docker compose down -v --remove-orphans

migrate_up:
	migrate -path internal/server/infrastructure/mysql/db/schema -database "mysql://root:$(MYSQL_ROOT_PASSWORD)@tcp(localhost:3306)/${MYSQL_DATABASE}" up

migrate_down:
	migrate -path internal/server/infrastructure/mysql/db/schema -database "mysql://root:$(MYSQL_ROOT_PASSWORD)@tcp(localhost:3306)/${MYSQL_DATABASE}" down

exec_db:
	docker compose exec db mysql -u root -p$(MYSQL_ROOT_PASSWORD) ${MYSQL_DATABASE}

test:
	go test -cover ./... -gcflags="all=-N -l" -v -coverprofile=cover.out
	go tool cover -html=cover.out

generate: docker-generate-mock
	docker run --rm -v $(PWD):/app ${GENERATE_IMAGE} sh -c "sqlc generate"

mock-gen: docker-generate-mock
	docker run --rm -v $(PWD):/app ${GENERATE_IMAGE} sh -c \
	"mockgen -package domain -source=internal/server/domain/log_repository.go -destination=internal/server/domain/log_mock.go && \
	mockgen -package usecase -source=internal/server/usecase/insert_log.go -destination=internal/server/usecase/insert_log_mock.go \
	mockgen -package usecase -source=internal/server/usecase/list_log.go -destination=internal/server/usecase/list_log_mock.go"

docker-generate-mock:
	docker build -f Dockerfile.generate -t ${GENERATE_IMAGE} .

send_message:
	RABBITMQ_URL=amqp://guest:guest@localhost:5672/ go run ./cmd/client $(ARGS)

create_doc:
	godoc -http=localhost:${GODOC_PORT}

format:
	go fmt ./...