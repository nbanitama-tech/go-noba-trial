package httpapi

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/bytedance/go-noba-trial/internal/domain"
)

type Middleware func(http.Handler) http.Handler

func chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}

func bearerAuthMiddleware(token string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authToken, ok := parseBearerToken(r.Header.Get("Authorization"))
			if !ok || subtle.ConstantTimeCompare([]byte(authToken), []byte(token)) != 1 {
				writeJSON(w, http.StatusUnauthorized, domain.GeneralResponse[any]{
					Success: false,
					Message: "unauthorized",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func parseBearerToken(header string) (string, bool) {
	const prefix = "Bearer "

	if !strings.HasPrefix(header, prefix) {
		return "", false
	}

	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	return token, token != ""
}
