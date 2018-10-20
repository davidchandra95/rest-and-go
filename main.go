package main

import (
	"log"
	"net/http"
	"os"
	"github.com/gorilla/handlers"
	"rest-and-go/store"
)

func main() {
	// Get the PORT env variable
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := store.NewRouter() // Create routes

	// Important lines for frond-end utilise
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})

	log.Fatal(http.ListenAndServe(":" + port, handlers.CORS(allowedOrigins, allowedMethods)(router)))
}