package httpapi

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bytedance/go-noba-trial/internal/domain"
	"github.com/bytedance/go-noba-trial/internal/usecase"
)

func TestPing(t *testing.T) {
	router := NewRouter(RouterConfig{
		HealthUsecase: usecase.NewHealthUsecase("test-service"),
		Logger:        slog.New(slog.NewTextHandler(io.Discard, nil)),
	})

	request := httptest.NewRequest(http.MethodGet, "/ping", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var body domain.GeneralResponse[domain.HealthStatus]
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !body.Success {
		t.Fatal("expected success response")
	}

	if body.Message != "pong" {
		t.Fatalf("expected message pong, got %q", body.Message)
	}

	if body.Data.Service != "test-service" {
		t.Fatalf("expected service test-service, got %q", body.Data.Service)
	}

	if body.Data.Status != "ok" {
		t.Fatalf("expected status ok, got %q", body.Data.Status)
	}
}

func TestPingRejectsUnsupportedMethod(t *testing.T) {
	router := NewRouter(RouterConfig{
		HealthUsecase: usecase.NewHealthUsecase("test-service"),
		Logger:        slog.New(slog.NewTextHandler(io.Discard, nil)),
	})

	request := httptest.NewRequest(http.MethodPost, "/ping", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, response.Code)
	}
}
