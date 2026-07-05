# go-noba-trial

Simple Go REST service using clean architecture boundaries.

## Structure

- `cmd/api`: application entrypoint
- `internal/config`: configuration loading
- `internal/datastore`: PostgreSQL connection and repositories
- `internal/domain`: domain models and response contracts
- `internal/usecase`: application use cases
- `internal/delivery/http`: REST handlers and routing

## Run

```sh
go run ./cmd/api
```

The service listens on `:8080` by default. Override it with:

```sh
HTTP_ADDRESS=:3000 go run ./cmd/api
```

## Docker

Build the image:

```sh
docker build -t go-noba-trial .
```

Run the container with an existing PostgreSQL database:

```sh
docker run --rm -p 8080:8080 \
  -e DATABASE_URL='postgres://noba:noba_password@host.docker.internal:5432/noba?sslmode=disable' \
  go-noba-trial
```

Run with a custom service name:

```sh
docker run --rm -p 8080:8080 \
  -e SERVICE_NAME=my-service \
  -e DATABASE_URL='postgres://noba:noba_password@host.docker.internal:5432/noba?sslmode=disable' \
  go-noba-trial
```

## Docker Compose

Start the API with PostgreSQL:

```sh
docker compose up --build
```

If port `8080` is already in use on your machine:

```sh
API_PORT=18080 docker compose up --build
```

The compose stack creates PostgreSQL with authenticated access:

- database: `noba`
- user: `noba`
- password: `noba_password`
- API database URL: `postgres://noba:noba_password@postgres:5432/noba?sslmode=disable`

## Endpoint

```sh
curl http://localhost:8080/ping
```

Response:

```json
{
  "success": true,
  "message": "pong",
  "data": {
    "service": "go-noba-trial",
    "status": "ok"
  }
}
```

Add a user:

```sh
curl -X POST http://localhost:8080/add \
  -H 'Content-Type: application/json' \
  -d '{
    "fullname": "Jane Doe",
    "email": "jane@example.com",
    "description": "Example user"
  }'
```

Response:

```json
{
  "success": true,
  "message": "user added successfully",
  "data": {
    "uuid": "generated-user-uuid",
    "fullname": "Jane Doe",
    "email": "jane@example.com",
    "description": "Example user"
  }
}
```

List users:

```sh
curl http://localhost:8080/user/list
```

Response:

```json
{
  "success": true,
  "message": "users fetched successfully",
  "data": [
    {
      "uuid": "f4b2fe41-4b68-42a9-8db2-8563dc5c7eb9",
      "fullname": "Jane Doe",
      "email": "jane@example.com",
      "description": "Example user"
    }
  ]
}
```

## Test

```sh
go test ./...
```
