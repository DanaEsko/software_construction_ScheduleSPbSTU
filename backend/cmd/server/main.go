// Command server starts the ScheduleSPbSTU backend.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/DanaEsko/software_construction_ScheduleSPbSTU/backend/internal/schedule"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	addr := getEnv("ADDR", ":8080")
	ruzURL := getEnv("RUZ_BASE_URL", "https://ruz.spbstu.ru/api/v1/ruz")

	svc := schedule.NewService(
		schedule.NewMemoryStore(),
		schedule.NewRuzClient(ruzURL, 3*time.Second),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /schedule", scheduleHandler(svc))

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("server started", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", "error", err)
			stop()
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	}
	slog.Info("server stopped")
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
		slog.Error("write health response", "error", err)
	}
}

func scheduleHandler(svc *schedule.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupID, err := strconv.Atoi(r.URL.Query().Get("group_id"))
		if err != nil || groupID <= 0 {
			writeError(w, http.StatusBadRequest, "validation_failed", "group_id must be a positive integer")
			return
		}

		date, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
		if err != nil {
			writeError(w, http.StatusBadRequest, "validation_failed", "date must be in YYYY-MM-DD format")
			return
		}

		res, err := svc.Get(r.Context(), groupID, schedule.WeekStart(date))
		if err != nil {
			slog.Error("get schedule", "group_id", groupID, "error", err)
			writeError(w, http.StatusServiceUnavailable, "schedule_unavailable", "schedule is temporarily unavailable")
			return
		}

		slog.Info("schedule served", "group_id", groupID, "state", res.State, "from_cache", res.FromCache)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cache-State", string(res.State))
		if _, err := w.Write(res.Payload); err != nil {
			slog.Error("write schedule response", "error", err)
		}
	}
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	body := map[string]string{"error": code, "message": message}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Error("write error response", "error", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
