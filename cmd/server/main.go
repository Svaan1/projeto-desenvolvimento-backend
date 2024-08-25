package main

import (
	"log"
	"net/http"

	"backendProject/routes"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}
}

func main() {
	r := routes.NewRouter()

	log.Print("listening :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
