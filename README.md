# go-noba-trial

Simple Go REST service using clean architecture boundaries.

## Structure

- `cmd/api`: application entrypoint
- `internal/config`: configuration loading
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

Run the container:

```sh
docker run --rm -p 8080:8080 go-noba-trial
```

Run with a custom service name:

```sh
docker run --rm -p 8080:8080 -e SERVICE_NAME=my-service go-noba-trial
```

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

## Test

```sh
go test ./...
```
