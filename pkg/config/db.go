package config

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
 "os"
	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	portStr := os.Getenv("DB_PORT")
	portInt, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port: %v", err)
	}

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		os.Getenv("DB_HOST"),
		portInt,
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Connection error:", err)
	}
	return db
}


