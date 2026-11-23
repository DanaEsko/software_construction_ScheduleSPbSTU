package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)
	connectDB()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /auth/register", registerHandler)
	mux.HandleFunc("POST /auth/login", loginHandler)

	port := getEnv("PORT", "8080")
	handler := loggingMiddleware(corsMiddleware(mux))

	log.Printf("server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
