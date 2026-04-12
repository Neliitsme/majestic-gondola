# majestic-gondola
A project to learn how to wield the power of golang.

## Topic
Some service that works with data related to music. Nothing too exciting. Like names, genres, authors, etc.

## Plans
- [ ] A simple CRUD with a few entities related to each other
- [ ] Write some business logic
- [ ] Communicate with other APIs ?
- [ ] Break down the app into microservices ?
- [ ] (Impossible) Fix the perpetual TODOs

## Usage

To run the app:

1. Start the database with
```sh
docker compose up -d
```

2. cd into the migartion folder and create the database schema with
```sh
go run . usage # check the list of available commands
go run . init
go run .
```

3. Start the app from the project's root with
```sh
go run ./cmd/server
```

(Optionally) generate the swag docs with
```sh
swag init -g cmd/server/main.go --parseDependency --dir ./,./internal/handlers,./internal/models
```