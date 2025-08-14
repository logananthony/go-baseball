package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func getenv(key, def string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	return v
}

func ConnectDB() *sql.DB {
	host := getenv("DB_HOST", "127.0.0.1")
	port := getenv("DB_PORT", "5432") // keep string
	user := getenv("DB_USER", "postgres")
	pass := getenv("DB_PASSWORD", "")
	name := getenv("DB_NAME", "postgres")
	// Default to TLS in prod; override locally with DB_SSLMODE=disable if needed.
	ssl := getenv("DB_SSLMODE", "require")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, pass, name, ssl,
	)

	log.Printf("Connecting to Postgres host=%s port=%s dbname=%s sslmode=%s", host, port, name, ssl)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		// This is a config error (bad DSN/driver); fail fast.
		log.Fatalf("sql.Open: %v", err)
	}

	// Start a background ping with backoff so the web server can start
	// and pass App Platform's readiness probe even if the DB isn't ready yet.
	go func() {
		backoff := time.Second
		for {
			if err := db.Ping(); err != nil {
				log.Printf("db.Ping failed: %v (retrying in %v)", err, backoff)
				time.Sleep(backoff)
				if backoff < 30*time.Second {
					backoff *= 2
				}
				continue
			}
			log.Printf("db.Ping succeeded; database connection established.")
			return
		}
	}()

	return db
}
