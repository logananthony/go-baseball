// package config

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"strconv"
//  "os"
// 	_ "github.com/lib/pq"
// )

// func ConnectDB() *sql.DB {
// 	portStr := os.Getenv("DB_PORT")
// 	portInt, err := strconv.Atoi(portStr)
// 	if err != nil {
// 		log.Fatalf("Invalid port: %v", err)
// 	}

// 	connStr := fmt.Sprintf(
// 		"host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
// 		os.Getenv("DB_HOST"),
// 		portInt,
// 		os.Getenv("DB_USER"),
// 		os.Getenv("DB_PASSWORD"),
// 		os.Getenv("DB_NAME"),
// 	)

// 	db, err := sql.Open("postgres", connStr)
// 	if err != nil {
// 		log.Fatal("Connection error:", err)
// 	}
// 	return db
// }

package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

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
	port := getenv("DB_PORT", "5432") // ← keep as string, no Atoi
	user := getenv("DB_USER", "postgres")
	pass := getenv("DB_PASSWORD", "")
	name := getenv("DB_NAME", "postgres")
	ssl := getenv("DB_SSLMODE", "disable") // local default; override to "require" for RDS

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, pass, name, ssl,
	)

	// Optional: redact password in logs if you log DSN
	log.Printf("Connecting to Postgres host=%s port=%s dbname=%s sslmode=%s", host, port, name, ssl)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("sql.Open: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("db.Ping: %v", err)
	}
	return db
}
