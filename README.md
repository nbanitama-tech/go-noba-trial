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

The service reads settings from `config.yaml` by default:

```yaml
http_address: ":8080"
service_name: "go-noba-trial"
database_url: "postgres://noba:noba_password@localhost:5432/noba?sslmode=disable"
bearer_token: "go-noba-trial-token"
```

Use another config file with `CONFIG_PATH`:

```sh
CONFIG_PATH=./config.local.yaml go run ./cmd/api
```

## Docker

Build the image:

```sh
docker build -t go-noba-trial .
```

Run the container with an existing PostgreSQL database by mounting a config file whose `database_url` points at that database:

```sh
docker run --rm -p 8080:8080 \
  -v "$PWD/config.local.yaml:/app/config.yaml:ro" \
  go-noba-trial
```

Run with a custom config file:

```sh
docker run --rm -p 8080:8080 \
  -v "$PWD/config.local.yaml:/app/config.yaml:ro" \
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
- API config file: `config.docker.yaml`

## Endpoint

```sh
curl -H 'Authorization: Bearer go-noba-trial-token' \
  http://localhost:8080/ping
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
  -H 'Authorization: Bearer go-noba-trial-token' \
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
curl -H 'Authorization: Bearer go-noba-trial-token' \
  http://localhost:8080/user/list
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
