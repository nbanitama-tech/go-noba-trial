package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	defaultConfigPath = "config.yaml"
)

type Config struct {
	HTTPAddress string `yaml:"http_address"`
	ServiceName string `yaml:"service_name"`
	DatabaseURL string `yaml:"database_url"`
	BearerToken string `yaml:"bearer_token"`
}

func Load() (Config, error) {
	return LoadFile(getEnv("CONFIG_PATH", defaultConfigPath))
}

func LoadFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (cfg Config) Validate() error {
	required := map[string]string{
		"http_address": cfg.HTTPAddress,
		"service_name": cfg.ServiceName,
		"database_url": cfg.DatabaseURL,
		"bearer_token": cfg.BearerToken,
	}

	for name, value := range required {
		if value == "" {
			return fmt.Errorf("missing required config field %q", name)
		}
	}

	return nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
