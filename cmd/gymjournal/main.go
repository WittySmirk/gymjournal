package main

import (
	"log"
	"net/http"
	"os"

	"github.com/WittySmirk/gymjournal/internal/auth"
	"github.com/WittySmirk/gymjournal/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("ENVIRONMENT")
	if env == "" || env == "development" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found")
		}
	}

	server.CreateDb()
	defer server.CloseDb()

	auth.CreateAuth()
	mux := server.CreateRoutes()

	http.ListenAndServe(":8080", mux)
}
