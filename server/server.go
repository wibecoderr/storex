package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/wibecoderr/storex/handler"
	"github.com/wibecoderr/storex/middleware"
)

func SetUpRoutes(r chi.Router) {
	r.Post("/register", handler.RegisterUser)
	r.Post("/login", handler.LoginUser)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Post("/logout", handler.LogoutUser)
		r.Get("/assets", handler.DisplayAsset)
		r.Get("/assets/{id}", handler.GetAssetByID)
		r.Get("/employees/asset", handler.ListAssetsByEmployee)
		r.Get("/employees", handler.GetEmpoloyee)

	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Use(middleware.RoleMiddleware("admin"))
		r.Post("/register/employee", handler.CreateEmployee)
		r.Post("/assets", handler.CreateAsset)
		r.Put("/assets/update/{id}", handler.UpdateAsset)
		r.Get("/assets/employee{id}", handler.ListAssetsByEmployeeAdmin)
		r.Post("/assets/return/{id}", handler.ReturnAssest)
		r.Post("/assets/assign", handler.AssignAsset)
		r.Delete("/user/{id}", handler.ArchieveUser)
		r.Delete("/assets/{id}", handler.DeleteAsset)
	})
}
