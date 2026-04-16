MAIN_PATH=./cmd/server/main.go

.PHONY: run swag setup

run: swag
	go run $(MAIN_PATH)

swag:
	swag init -g $(MAIN_PATH) --parseDependency --dir ./,./internal/handlers,./internal/models

setup:
	cp .env.example .env