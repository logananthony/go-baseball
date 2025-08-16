package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/logananthony/go-baseball/pkg/api"
	"github.com/logananthony/go-baseball/pkg/config"
)

func init() {
	// Loads .env when present (dev). In App Platform, set env vars in the UI/spec.
	_ = godotenv.Load()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	db := config.ConnectDB()
	// Reasonable pool defaults; tweak to taste.
	db.SetMaxOpenConns(15)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := api.NewAPIServer(":"+port, db)
	if err := server.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
