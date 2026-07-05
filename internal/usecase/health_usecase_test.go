package usecase

import "testing"

func TestHealthUsecasePing(t *testing.T) {
	usecase := NewHealthUsecase("test-service")

	status := usecase.Ping()

	if status.Service != "test-service" {
		t.Fatalf("expected service test-service, got %q", status.Service)
	}

	if status.Status != "ok" {
		t.Fatalf("expected status ok, got %q", status.Status)
	}
}
