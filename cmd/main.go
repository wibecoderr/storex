package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/wibecoderr/storex/database"
	"github.com/wibecoderr/storex/handler"
	"github.com/wibecoderr/storex/middleware"
)

func main() {
	err := database.ConnectAndMigrate(
		"localhost",
		"5433",
		"postgres",
		"local",
		"local",
		database.SSLModeDisable,
	)

	// trasncation , storex -> utils ,
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer database.ShutdownDatabase()

	r := chi.NewRouter()

	r.Post("/register", handler.RegisterUser)
	r.Post("/login", handler.LoginUser)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/assets", handler.DisplayAsset) // required p
		r.Get("/assets/{id}", handler.GetAssetByID)
		r.Get("/employees/{id}/assets", handler.ListAssetsByEmployee)
		r.Post("/logout", handler.LogoutUser)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Use(middleware.RoleMiddleware("admin"))
		r.Post("/assets", handler.CreateAsset)
		r.Post("/assets/assign", handler.AssignAsset)
		r.Delete("/assets/{id}", handler.DeleteAsset)
	})

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
