package config

import (
	"os"
	"testing"
)

func TestLoadUsesConfiguredPath(t *testing.T) {
	configPath := writeConfig(t, `
http_address: ":9090"
service_name: "custom-service"
database_url: "postgres://example"
bearer_token: "custom-token"
`)
	t.Setenv("CONFIG_PATH", configPath)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected config to load, got error %v", err)
	}

	if cfg.HTTPAddress != ":9090" {
		t.Fatalf("expected HTTP address from config file, got %q", cfg.HTTPAddress)
	}

	if cfg.ServiceName != "custom-service" {
		t.Fatalf("expected service name from config file, got %q", cfg.ServiceName)
	}

	if cfg.DatabaseURL != "postgres://example" {
		t.Fatalf("expected database URL from config file, got %q", cfg.DatabaseURL)
	}

	if cfg.BearerToken != "custom-token" {
		t.Fatalf("expected bearer token from config file, got %q", cfg.BearerToken)
	}
}

func TestLoadFileReturnsErrorForInvalidYAML(t *testing.T) {
	configPath := writeConfig(t, "http_address: [")

	if _, err := LoadFile(configPath); err == nil {
		t.Fatal("expected invalid config file to return an error")
	}
}

func TestLoadFileReturnsErrorForMissingRequiredField(t *testing.T) {
	configPath := writeConfig(t, `
http_address: ":9090"
service_name: "custom-service"
database_url: "postgres://example"
`)

	if _, err := LoadFile(configPath); err == nil {
		t.Fatal("expected missing config field to return an error")
	}
}

func writeConfig(t *testing.T, content string) string {
	t.Helper()

	configPath := t.TempDir() + "/config.yaml"
	if err := os.WriteFile(configPath, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	return configPath
}
