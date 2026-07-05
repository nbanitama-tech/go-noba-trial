package config

import "os"

const (
	defaultHTTPAddress = ":8080"
	defaultServiceName = "go-noba-trial"
)

type Config struct {
	HTTPAddress string
	ServiceName string
}

func Load() Config {
	return Config{
		HTTPAddress: getEnv("HTTP_ADDRESS", defaultHTTPAddress),
		ServiceName: getEnv("SERVICE_NAME", defaultServiceName),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
