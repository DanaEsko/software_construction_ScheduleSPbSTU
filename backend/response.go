package main

import (
	"encoding/json"
	"net/http"
)

// отправляет ответ в формате JSON с нужным кодом
func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
