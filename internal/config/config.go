package config

import "os"

const (
	defaultHTTPAddress = ":8080"
	defaultServiceName = "go-noba-trial"
	defaultDatabaseURL = "postgres://noba:noba_password@localhost:5432/noba?sslmode=disable"
)

type Config struct {
	HTTPAddress string
	ServiceName string
	DatabaseURL string
}

func Load() Config {
	return Config{
		HTTPAddress: getEnv("HTTP_ADDRESS", defaultHTTPAddress),
		ServiceName: getEnv("SERVICE_NAME", defaultServiceName),
		DatabaseURL: getEnv("DATABASE_URL", defaultDatabaseURL),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
