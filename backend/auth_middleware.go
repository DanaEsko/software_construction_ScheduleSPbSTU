package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const userIDKey ctxKey = "userID"

// проверяет токен из заголовка Authorization: Bearer <токен>
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			writeError(w, http.StatusUnauthorized, "unauthorized", "missing or invalid token")
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			return jwtSecret(), nil
		}, jwt.WithValidMethods([]string{"HS256"}))
		if err != nil || !token.Valid {
			writeError(w, http.StatusUnauthorized, "unauthorized", "missing or invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			writeError(w, http.StatusUnauthorized, "unauthorized", "missing or invalid token")
			return
		}
		userID, _ := claims["sub"].(string)
		if userID == "" {
			writeError(w, http.StatusUnauthorized, "unauthorized", "missing or invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}

// достаёт id пользователя из контекста
func currentUserID(r *http.Request) string {
	id, _ := r.Context().Value(userIDKey).(string)
	return id
}
