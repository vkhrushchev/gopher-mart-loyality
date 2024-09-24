package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

var log = zap.Must(zap.NewDevelopment()).Sugar()

const AuthTokenCoockieName = "AuthToken"

type MidlewareContextKey string

const (
	UserLoginContextKey MidlewareContextKey = "user_login"
)

func NewJWTAuthMiddleware(jwtSecretKey string) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authTokenCookie, err := r.Cookie(AuthTokenCoockieName)
			if errors.Is(err, http.ErrNoCookie) {
				log.Infow("middleware: no auth token cookie found")

				w.WriteHeader(http.StatusUnauthorized)
				return
			} else if err != nil {
				log.Errorw("middleware: unexpected error", "error", err.Error())

				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			authToken := authTokenCookie.Value
			jwtClaims := jwt.MapClaims{}
			_, err = jwt.ParseWithClaims(
				authToken,
				jwtClaims,
				func(t *jwt.Token) (interface{}, error) {
					return []byte(jwtSecretKey), nil
				},
			)
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				log.Infow("middleware: jwt signature is invalid")

				w.WriteHeader(http.StatusUnauthorized)
				return
			} else if err != nil {
				log.Errorw("middleware: unexpected error", "error", err.Error())

				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), UserLoginContextKey, jwtClaims["sub"]))

			handler.ServeHTTP(w, r)
		})
	}
}

func CheckContentTypeMiddleware(expectedContentType string, next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if contentType != expectedContentType {
			log.Infow(
				"middleaware: not supported \"Content-Type\" header",
				"Content-Type", contentType,
			)

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		next(w, r)
	}
}
