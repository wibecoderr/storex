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

	// history dashboard logout return assest  -- api , enum in go ,
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer database.ShutdownDatabase()

	r := chi.NewRouter()

	r.Post("/register", handler.RegisterUser) // workign correctly
	r.Post("/login", handler.LoginUser)       // working

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Post("/logout", handler.LogoutUser)  // working
		r.Get("/assets", handler.DisplayAsset) // required p
		r.Get("/assets/{id}", handler.GetAssetByID)
		r.Get("/employees", handler.ListAssetsByEmployee) //working
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Use(middleware.RoleMiddleware("admin"))
		r.Post("/register/employee", handler.CreateEmployee) //working
		r.Post("/assets", handler.CreateAsset)               //working
		r.Get("/assets/{id}", handler.ListAssetsByEmployeeAdmin)
		r.Post("/assets/return/{id}", handler.ReturnAssest) // working
		r.Post("/assets/assign", handler.AssignAsset)       // working
		r.Delete("/assets/{id}", handler.DeleteAsset)       // working
	})

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
