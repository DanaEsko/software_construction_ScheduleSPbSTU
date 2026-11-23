package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func writeError(w http.ResponseWriter, code int, errCode, message string) {
	writeJSON(w, code, map[string]string{"error": errCode, "message": message})
}

// POST /auth/register
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		writeError(w, http.StatusServiceUnavailable, "db_unavailable", "database is not connected")
		return
	}

	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "validation_failed", "invalid JSON body")
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	if !strings.Contains(email, "@") {
		writeError(w, http.StatusBadRequest, "validation_failed", "email is invalid")
		return
	}
	if len(req.Password) < 8 {
		writeError(w, http.StatusBadRequest, "validation_failed", "password must be at least 8 characters")
		return
	}

	hash, err := hashPassword(req.Password)
	if err != nil {
		log.Printf("hash error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "internal server error, please try again later")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var userID string
	err = db.QueryRow(ctx,
		`INSERT INTO users (id, email, password_hash)
		 VALUES (gen_random_uuid(), $1, $2)
		 RETURNING id::text`,
		email, hash,
	).Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // email уже есть
			writeError(w, http.StatusConflict, "email_taken", "user with this email already exists")
			return
		}
		log.Printf("insert user error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "internal server error, please try again later")
		return
	}

	token, err := makeToken(userID)
	if err != nil {
		log.Printf("token error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "internal server error, please try again later")
		return
	}

	log.Printf("user registered: %s", userID)
	writeJSON(w, http.StatusCreated, map[string]string{"token": token})
}

// loginHandler: POST /auth/login
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		writeError(w, http.StatusServiceUnavailable, "db_unavailable", "database is not connected")
		return
	}

	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "validation_failed", "invalid JSON body")
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "validation_failed", "email and password are required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var userID, hash string
	err := db.QueryRow(ctx,
		`SELECT id::text, password_hash FROM users WHERE email = $1`,
		email,
	).Scan(&userID, &hash)

	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusUnauthorized, "invalid_credentials", "wrong email or password")
		return
	}
	if err != nil {
		log.Printf("select user error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "internal server error, please try again later")
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		writeError(w, http.StatusUnauthorized, "invalid_credentials", "wrong email or password")
		return
	}

	token, err := makeToken(userID)
	if err != nil {
		log.Printf("token error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "internal server error, please try again later")
		return
	}

	log.Printf("user logged in: %s", userID)
	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}
