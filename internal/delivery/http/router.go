package httpapi

import (
	"log/slog"
	"net/http"

	"github.com/bytedance/go-noba-trial/internal/usecase"
)

type RouterConfig struct {
	HealthUsecase usecase.HealthUsecase
	UserUsecase   usecase.UserUsecase
	Logger        *slog.Logger
}

func NewRouter(cfg RouterConfig) http.Handler {
	handler := NewHandler(cfg.HealthUsecase, cfg.UserUsecase, cfg.Logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", handler.Ping)
	mux.HandleFunc("/add", handler.AddUser)
	mux.HandleFunc("/user/list", handler.ListUsers)

	return mux
}
