package httpapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/bytedance/go-noba-trial/internal/domain"
	"github.com/bytedance/go-noba-trial/internal/usecase"
)

type Handler struct {
	healthUsecase usecase.HealthUsecase
	userUsecase   usecase.UserUsecase
	logger        *slog.Logger
}

func NewHandler(healthUsecase usecase.HealthUsecase, userUsecase usecase.UserUsecase, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}

	return &Handler{
		healthUsecase: healthUsecase,
		userUsecase:   userUsecase,
		logger:        logger,
	}
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, domain.GeneralResponse[any]{
			Success: false,
			Message: "method not allowed",
		})
		return
	}

	writeJSON(w, http.StatusOK, domain.GeneralResponse[domain.HealthStatus]{
		Success: true,
		Message: "pong",
		Data:    h.healthUsecase.Ping(),
	})
}

func (h *Handler) AddUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, domain.GeneralResponse[any]{
			Success: false,
			Message: "method not allowed",
		})
		return
	}

	var input domain.CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, domain.GeneralResponse[any]{
			Success: false,
			Message: "invalid json body",
		})
		return
	}

	user, err := h.userUsecase.Add(r.Context(), input)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidUserInput) {
			writeJSON(w, http.StatusBadRequest, domain.GeneralResponse[any]{
				Success: false,
				Message: "fullname and email are required",
			})
			return
		}

		h.logger.Error("failed to add user", "error", err)
		writeJSON(w, http.StatusInternalServerError, domain.GeneralResponse[any]{
			Success: false,
			Message: "failed to add user",
		})
		return
	}

	writeJSON(w, http.StatusCreated, domain.GeneralResponse[domain.User]{
		Success: true,
		Message: "user added successfully",
		Data:    user,
	})
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, domain.GeneralResponse[any]{
			Success: false,
			Message: "method not allowed",
		})
		return
	}

	users, err := h.userUsecase.List(r.Context())
	if err != nil {
		h.logger.Error("failed to list users", "error", err)
		writeJSON(w, http.StatusInternalServerError, domain.GeneralResponse[any]{
			Success: false,
			Message: "failed to list users",
		})
		return
	}

	total := len(users)
	if total == 0 {
		users = nil
	}

	writeJSON(w, http.StatusOK, domain.ListUsersResponse{
		Success: true,
		Message: "users fetched successfully",
		Data:    users,
		Total:   total,
	})
}

func writeJSON(w http.ResponseWriter, statusCode int, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("failed to write json response", "error", err)
	}
}
