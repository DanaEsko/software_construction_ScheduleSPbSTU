package main

import "os"

// возвращает переменную окружения или значение по умолчанию
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
