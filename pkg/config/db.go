package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		Host, Port, User, Password, Dbname,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Connection error:", err)
	}
	return db
}

