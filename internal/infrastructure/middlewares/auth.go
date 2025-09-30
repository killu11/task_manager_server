package middlewares

import (
	"context"
	"net/http"
	"strings"
	"task_manager_server/pkg/security"
)

func AuthJWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}
		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "missing token from auth-header", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]
		valid, err := security.CheckTokenValid(tokenString)

		if err != nil || !valid {
			http.Error(w, "invalid token, sign-in or sign-up again", http.StatusUnauthorized)
			return
		}
		id, err := security.ParseUserID(tokenString)

		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", id)
		reqWithID := r.WithContext(ctx)
		next.ServeHTTP(w, reqWithID)
	})
}
