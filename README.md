# Implementing Change: Business Simulation

Currently this is about playing with [Fiber](https://gofiber.io) framework in Go.

It's an API that can:
- `POST /login` to get access and refresh tokens that are cached to Redis
- `POST /v1/logout` to empty currently cached acceess token in Redis
- `GET /v1/ping` is a dummy page that checks if the correct token is provided
- Users are persisted in Postgres DB
- Redis is used for caching tokens

## TODO
- Persist users with Postgres
- Docker & Docker Compose
- Tests
- Add Github Action for tests
- Add Sentry for catching errors

## How to run

- Have Go 1.15+ installed
- Have Redis installed and running on local host without password
- Have Postgres 12 installed and running
- Apply DB migration SQL script from `migrations/migration.sql` using any Postgres query tool.
- Have `.env` with the following parameters (set appropriate values if needed):
```
ACCESS_SECRET=notsosecret
REFRESH_SECRET=alsonotverysecret
REDIS_DSN=localhost:6379
POSTGRES_URL=postgresql://localhost:5432/change
```  
- `$ go get`
- `$ go run .`
- API is available at `http://127.0.0.1:3000`