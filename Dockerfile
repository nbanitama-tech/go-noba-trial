FROM golang:1.24.1-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api ./cmd/api

FROM alpine:3.21

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /bin/api /app/api

USER app

EXPOSE 8080

ENV HTTP_ADDRESS=:8080

ENTRYPOINT ["/app/api"]
