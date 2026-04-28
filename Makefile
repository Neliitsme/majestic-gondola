include .env

MAIN_PATH=./cmd/server/main.go
BINARY_NAME=server
MIGRATIONS_PATH=./db/migrations

.PHONY: *

run: swag
	go run $(MAIN_PATH)

swag:
	swag fmt
	swag init -g $(MAIN_PATH) --parseDependency --dir ./,./internal/handlers --useStructName -q

setup:
	cp .env.example .env

build:
	go build -o $(BINARY_NAME) $(MAIN_PATH)

clean:
	rm -f $(BINARY_NAME)

test:
	go test -v ./...

mup:
	migrate -database ${POSTGRES_URL} -path ${MIGRATIONS_PATH} up

mdown:
	migrate -database ${POSTGRES_URL} -path ${MIGRATIONS_PATH} down