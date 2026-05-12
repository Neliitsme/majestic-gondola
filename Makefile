include .env

MAIN_PATH=./cmd/server/main.go
MAIN_PKG=./cmd/server
MIGRATIONS_PATH=./db/migrations

.PHONY: run swag setup build clean test mup mdown seed

run: swag
	go run $(MAIN_PKG)

swag:
	swag fmt
	swag init -g $(MAIN_PATH) --parseDependency --dir ./,./internal/handlers --useStructName -q

setup:
	cp .env.example .env

build:
	go build $(MAIN_PKG)

test:
	go test -v ./...

mup:
	migrate -database ${POSTGRES_URL} -path ${MIGRATIONS_PATH} up

mdown:
	migrate -database ${POSTGRES_URL} -path ${MIGRATIONS_PATH} down

seed:
	psql ${POSTGRES_URL} < ./db/seed.sql
