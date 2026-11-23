package main

import (
	"log"
	"net/http"
	"os"
)

<<<<<<< HEAD
=======
var db *pgxpool.Pool // nil, пока база не подключена

// Логи
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s (%s)", r.Method, r.URL.Path, time.Since(start))
	})
}

// CORS
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

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

// /health с проверкой базы 
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

// Подключение к PostgreSQL
func connectDB() {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		log.Println("DATABASE_URL is not set, starting without database")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		log.Printf("database config error: %v", err)
		return
	}
	if err := pool.Ping(ctx); err != nil {
		log.Printf("database connection failed: %v", err)
	} else {
		log.Println("database connection successful")
	}
	db = pool
}

>>>>>>> 4ae548c447dde5269b77f8a6bc1b4ffca4f334fc
func main() {
	log.SetOutput(os.Stdout)

	connectDB()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)

	port := getEnv("PORT", "8080")
	handler := loggingMiddleware(corsMiddleware(mux))

	log.Printf("server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
