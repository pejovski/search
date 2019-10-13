package http

import (
	"net/http"
	"os"
	"strings"
)

func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		apiKey := strings.TrimSpace(r.Header.Get("X-API-Key"))

		if apiKey != os.Getenv("APP_KEY") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
