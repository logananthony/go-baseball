package main

import (
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

	db := config.ConnectDB()
	defer db.Close()

	server := api.NewAPIServer(":8080", db)
	if err := server.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
