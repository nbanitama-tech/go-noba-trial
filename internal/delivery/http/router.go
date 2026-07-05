package httpapi

import (
	"log/slog"
	"net/http"

	"github.com/bytedance/go-noba-trial/internal/usecase"
)

type RouterConfig struct {
	HealthUsecase usecase.HealthUsecase
	Logger        *slog.Logger
}

func NewRouter(cfg RouterConfig) http.Handler {
	handler := NewHandler(cfg.HealthUsecase, cfg.Logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", handler.Ping)

	return mux
}
