package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"backendProject/internal/db"
	"backendProject/routes"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("error loading .env file")
	}
}

func main() {
	ctx := context.Background()

	rdb, err := db.NewRedisDB(ctx)
	if err != nil {
		log.Fatalf("error connecting to redis: %v", err)
	}

	r := routes.NewRouter(rdb)
	port := os.Getenv("SERVER_PORT")
	log.Printf("listening :%s", port)
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
