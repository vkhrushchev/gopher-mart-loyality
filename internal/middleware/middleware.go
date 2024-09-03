package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

var log = zap.Must(zap.NewDevelopment()).Sugar()

func JWTAuth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infow("JWTAuth middleware called.", "path", r.URL.Path)
		w.WriteHeader(http.StatusUnauthorized)

		// handler.ServeHTTP(w, r)
	})
}
