# Implementing Change: Business Simulation

Currently this is about playing with [Fiber](https://gofiber.io) framework in Go.

It's an API that can:
- `POST /login` to get access and refresh tokens that are cached to Redis
- `POST /v1/logout` to empty currently cached acceess token in Redis
- `GET /v1/ping` is a dummy page that checks if the correct token is provided

## TODO
- Persist users with Postgres
- Docker & Docker Compose
- Tests
- Add Github Action for tests
- Add Sentry for catching errors

## How to run

- Have Go 1.15+ installed
- Have Redis installed and running on local host without password
- `$ go get`
- `$ go run .`
- API is available at `http://127.0.0.1:3000`