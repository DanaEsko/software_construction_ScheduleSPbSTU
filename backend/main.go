package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

// Логи: пишет в консоль каждый запрос (подзадача 8)
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s (%s)", r.Method, r.URL.Path, time.Since(start))
	})
}

// CORS: разрешает запросы от фронтенда (подзадача 9)
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func main() {
	log.SetOutput(os.Stdout)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)

	handler := loggingMiddleware(corsMiddleware(mux))

	log.Println("server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
