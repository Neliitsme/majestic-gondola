# majestic-gondola
A project to learn how to wield the power of golang.

## Topic
Some service that works with data related to music. Nothing too exciting. Like names, genres, authors, etc.

## Plans
- [x] A simple CRUD with a few entities related to each other
- [ ] Write some business logic
- [ ] Write some frontend
- [ ] Add in-memory cache
- [ ] Try zenrpc
- [ ] Communicate with other APIs ?
- [ ] Break down the app into microservices ?
- [ ] (Impossible) Fix the perpetual TODOs

## Dependencies/Tech

- go 1.26.2
- Docker (postgres:18.3-bookworm)
- [golang-migrate](https://github.com/golang-migrate/migrate) for managing migrations

## Usage

To run the app:

1. Start the database with
```sh
docker compose up -d
```

2. Set up the `.env` file by copying the `.env.example` with the following command.

*Unless* you modify anything credential-sensitive, the example file has valid strings and keys to get this app up and running. It's a test project after all.
```sh
make init
```

3. Restore the database schema with
```sh
make mup
```

4. Start the app with
```sh
make
```

(Optionally) generate the swag docs with the following command. You only need to do it if you modify the API. Otherwise the up-to-date docs are present in the `/docs` folder and should be accessible by going to `<host>:<port>/swagger/index.html` when you run the app.
```sh
make swag
```