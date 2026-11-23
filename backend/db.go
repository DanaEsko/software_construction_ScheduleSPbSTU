package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool // nil, пока база не подключена

// подключается к PostgreSQL по строке из DATABASE_URL
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
