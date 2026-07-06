package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bytedance/go-noba-trial/internal/domain"
	"github.com/bytedance/go-noba-trial/internal/usecase"
)

const testBearerToken = "test-token"

type fakeUserUsecase struct {
	user       domain.User
	users      []domain.User
	err        error
	failOnCall bool
}

func (u fakeUserUsecase) Add(context.Context, domain.CreateUserInput) (domain.User, error) {
	if u.failOnCall {
		panic("user usecase should not be called")
	}

	return u.user, u.err
}

func (u fakeUserUsecase) List(context.Context) ([]domain.User, error) {
	if u.failOnCall {
		panic("user usecase should not be called")
	}

	return u.users, u.err
}

func newAuthorizedRequest(method string, target string, body io.Reader) *http.Request {
	request := httptest.NewRequest(method, target, body)
	request.Header.Set("Authorization", "Bearer "+testBearerToken)
	return request
}

func TestPing(t *testing.T) {
	router := NewRouter(RouterConfig{
		HealthUsecase: usecase.NewHealthUsecase("test-service"),
		Logger:        slog.New(slog.NewTextHandler(io.Discard, nil)),
		BearerToken:   testBearerToken,
	})

	request := newAuthorizedRequest(http.MethodGet, "/ping", nil)
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
		BearerToken:   testBearerToken,
	})

	request := newAuthorizedRequest(http.MethodPost, "/ping", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, response.Code)
	}
}

func TestAddUser(t *testing.T) {
	router := NewRouter(RouterConfig{
		HealthUsecase: usecase.NewHealthUsecase("test-service"),
		UserUsecase: fakeUserUsecase{
			user: domain.User{
				UUID:        "f4b2fe41-4b68-42a9-8db2-8563dc5c7eb9",
				Fullname:    "Jane Doe",
				Email:       "jane@example.com",
				Description: "Example user",
			},
		},
		Logger:      slog.New(slog.NewTextHandler(io.Discard, nil)),
		BearerToken: testBearerToken,
	})

	body := bytes.NewBufferString(`{"fullname":"Jane Doe","email":"jane@example.com","description":"Example user"}`)
	request := newAuthorizedRequest(http.MethodPost, "/add", body)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, response.Code)
	}

	var responseBody domain.GeneralResponse[domain.User]
	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !responseBody.Success {
		t.Fatal("expected success response")
	}

	if responseBody.Data.Email != "jane@example.com" {
		t.Fatalf("expected email jane@example.com, got %q", responseBody.Data.Email)
	}
}

func TestAddUserReturnsBadRequestForInvalidInput(t *testing.T) {
	router := NewRouter(RouterConfig{
		HealthUsecase: usecase.NewHealthUsecase("test-service"),
		UserUsecase: fakeUserUsecase{
			err: usecase.ErrInvalidUserInput,
		},
		Logger:      slog.New(slog.NewTextHandler(io.Discard, nil)),
		BearerToken: testBearerToken,
	})

	body := bytes.NewBufferString(`{"fullname":"","email":"","description":"Example user"}`)
	request := newAuthorizedRequest(http.MethodPost, "/add", body)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestAddUserReturnsInternalServerError(t *testing.T) {
	router := NewRouter(RouterConfig{
		HealthUsecase: usecase.NewHealthUsecase("test-service"),
		UserUsecase: fakeUserUsecase{
			err: errors.New("database failed"),
		},
		Logger:      slog.New(slog.NewTextHandler(io.Discard, nil)),
		BearerToken: testBearerToken,
	})

	body := bytes.NewBufferString(`{"fullname":"Jane Doe","email":"jane@example.com","description":"Example user"}`)
	request := newAuthorizedRequest(http.MethodPost, "/add", body)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, response.Code)
	}
}

func TestListUsers(t *testing.T) {
	router := NewRouter(RouterConfig{
		HealthUsecase: usecase.NewHealthUsecase("test-service"),
		UserUsecase: fakeUserUsecase{
			users: []domain.User{
				{
					UUID:        "f4b2fe41-4b68-42a9-8db2-8563dc5c7eb9",
					Fullname:    "Jane Doe",
					Email:       "jane@example.com",
					Description: "Example user",
				},
				{
					UUID:        "5c8a9be7-8d62-4d1f-8cf1-80cb5f1d3c55",
					Fullname:    "John Smith",
					Email:       "john@example.com",
					Description: "Another user",
				},
			},
		},
		Logger:      slog.New(slog.NewTextHandler(io.Discard, nil)),
		BearerToken: testBearerToken,
	})

	request := newAuthorizedRequest(http.MethodGet, "/user/list", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var responseBody domain.ListUsersResponse
	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !responseBody.Success {
		t.Fatal("expected success response")
	}

	if len(responseBody.Data) != 2 {
		t.Fatalf("expected 2 users, got %d", len(responseBody.Data))
	}

	if responseBody.Total != 2 {
		t.Fatalf("expected total 2, got %d", responseBody.Total)
	}

	if responseBody.Data[0].Email != "jane@example.com" {
		t.Fatalf("expected first email jane@example.com, got %q", responseBody.Data[0].Email)
	}
}

func TestListUsersReturnsNilDataWhenEmpty(t *testing.T) {
	router := NewRouter(RouterConfig{
		HealthUsecase: usecase.NewHealthUsecase("test-service"),
		UserUsecase:   fakeUserUsecase{},
		Logger:        slog.New(slog.NewTextHandler(io.Discard, nil)),
		BearerToken:   testBearerToken,
	})

	request := newAuthorizedRequest(http.MethodGet, "/user/list", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var responseBody domain.ListUsersResponse
	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !responseBody.Success {
		t.Fatal("expected success response")
	}

	if responseBody.Data != nil {
		t.Fatalf("expected nil data, got %#v", responseBody.Data)
	}

	if responseBody.Total != 0 {
		t.Fatalf("expected total 0, got %d", responseBody.Total)
	}
}

func TestRouterRejectsMissingBearerToken(t *testing.T) {
	router := NewRouter(RouterConfig{
		HealthUsecase: usecase.NewHealthUsecase("test-service"),
		UserUsecase: fakeUserUsecase{
			failOnCall: true,
		},
		Logger:      slog.New(slog.NewTextHandler(io.Discard, nil)),
		BearerToken: testBearerToken,
	})

	request := httptest.NewRequest(http.MethodGet, "/user/list", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, response.Code)
	}

	var responseBody domain.GeneralResponse[any]
	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if responseBody.Success {
		t.Fatal("expected failure response")
	}

	if responseBody.Message != "unauthorized" {
		t.Fatalf("expected unauthorized message, got %q", responseBody.Message)
	}
}

func TestRouterRejectsInvalidBearerToken(t *testing.T) {
	router := NewRouter(RouterConfig{
		HealthUsecase: usecase.NewHealthUsecase("test-service"),
		UserUsecase: fakeUserUsecase{
			failOnCall: true,
		},
		Logger:      slog.New(slog.NewTextHandler(io.Discard, nil)),
		BearerToken: testBearerToken,
	})

	request := httptest.NewRequest(http.MethodGet, "/user/list", nil)
	request.Header.Set("Authorization", "Bearer wrong-token")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, response.Code)
	}
}

func TestRouterRunsConfiguredMiddlewareAfterBearerAuth(t *testing.T) {
	router := NewRouter(RouterConfig{
		HealthUsecase: usecase.NewHealthUsecase("test-service"),
		Logger:        slog.New(slog.NewTextHandler(io.Discard, nil)),
		BearerToken:   testBearerToken,
		Middlewares: []Middleware{
			func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("X-Test-Middleware", "executed")
					next.ServeHTTP(w, r)
				})
			},
		},
	})

	request := newAuthorizedRequest(http.MethodGet, "/ping", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	if response.Header().Get("X-Test-Middleware") != "executed" {
		t.Fatal("expected configured middleware to run")
	}
}
