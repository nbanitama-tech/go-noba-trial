package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/bytedance/go-noba-trial/internal/domain"
	"github.com/bytedance/go-noba-trial/internal/usecase"
)

type Handler struct {
	healthUsecase usecase.HealthUsecase
	logger        *slog.Logger
}

func NewHandler(healthUsecase usecase.HealthUsecase, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}

	return &Handler{
		healthUsecase: healthUsecase,
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

func writeJSON(w http.ResponseWriter, statusCode int, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("failed to write json response", "error", err)
	}
}
