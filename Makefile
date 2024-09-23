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