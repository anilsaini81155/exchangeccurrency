package middleware

import (
	"net/http"
	"strings"

	"github.com/anilsaini81155/exchangeccurrency/utils"
)

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			next.ServeHTTP(w, r)
		} else {
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "Missing token", http.StatusUnauthorized)
				return
			}

			token = strings.TrimPrefix(token, "Bearer ")
			_, err := utils.ValidateJWT(token)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		}

	})
}
