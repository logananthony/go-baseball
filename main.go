package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/logananthony/go-baseball/pkg/api"
	"github.com/logananthony/go-baseball/pkg/config"
)

func init() {
	_ = godotenv.Load()
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func main() {

	defer timer("main")()
	time.Sleep(time.Second * 2)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Dbname)

	fmt.Println("DB HOST:", config.Host)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	server := api.NewAPIServer(":8080", db)
	if err := server.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
