package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Hello World!"}`)
	})

	return r
}

func main() {
	r := NewRouter()
	log.Print("listening :8080")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
