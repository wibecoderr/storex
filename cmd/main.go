package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/wibecoderr/storex/database"
	"github.com/wibecoderr/storex/server"
)

func main() {
	err := database.ConnectAndMigrate(
		/*
			DB_DATABASE=postgres;DB_HOST=localhost;DB_PASSWORD=local;DB_PORT=5433;DB_USER=local
		*/
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		database.SSLModeDisable,
	)

	// history dashboard logout return assest  -- api , enum in go ,
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer database.ShutdownDatabase()

	r := chi.NewRouter()
	server.SetUpRoutes(r)
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
