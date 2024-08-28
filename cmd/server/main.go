package main

import (
	"log"
	"net/http"
	"os"

	"backendProject/routes"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}
}

func main() {
	port := os.Getenv("SERVER_PORT")
	r := routes.NewRouter()

	log.Printf("listening :%s", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
