package config

import "testing"

func TestLoadUsesDefaults(t *testing.T) {
	t.Setenv("HTTP_ADDRESS", "")
	t.Setenv("SERVICE_NAME", "")
	t.Setenv("DATABASE_URL", "")

	cfg := Load()

	if cfg.HTTPAddress != defaultHTTPAddress {
		t.Fatalf("expected default HTTP address %q, got %q", defaultHTTPAddress, cfg.HTTPAddress)
	}

	if cfg.ServiceName != defaultServiceName {
		t.Fatalf("expected default service name %q, got %q", defaultServiceName, cfg.ServiceName)
	}

	if cfg.DatabaseURL != defaultDatabaseURL {
		t.Fatalf("expected default database URL %q, got %q", defaultDatabaseURL, cfg.DatabaseURL)
	}
}

func TestLoadUsesEnvironmentValues(t *testing.T) {
	t.Setenv("HTTP_ADDRESS", ":9090")
	t.Setenv("SERVICE_NAME", "custom-service")
	t.Setenv("DATABASE_URL", "postgres://example")

	cfg := Load()

	if cfg.HTTPAddress != ":9090" {
		t.Fatalf("expected HTTP address from environment, got %q", cfg.HTTPAddress)
	}

	if cfg.ServiceName != "custom-service" {
		t.Fatalf("expected service name from environment, got %q", cfg.ServiceName)
	}

	if cfg.DatabaseURL != "postgres://example" {
		t.Fatalf("expected database URL from environment, got %q", cfg.DatabaseURL)
	}
}
