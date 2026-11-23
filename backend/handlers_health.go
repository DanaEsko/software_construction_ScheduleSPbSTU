package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

// отвечает на GET /health и проверяет базу
func healthHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := db.Ping(ctx); err != nil {
		log.Printf("database ping failed: %v", err)
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"status":   "degraded",
			"database": "unavailable",
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"status":   "ok",
		"database": "ok",
	})
}
