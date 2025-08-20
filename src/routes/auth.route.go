package routes

import (
	"database/sql"
	"net/http"
	"restaurant-backend/src/controllers"
)

func AuthRoutes(server *http.ServeMux, db *sql.DB) {
	authController := controllers.NewAuthController(db)

	server.HandleFunc("/api/auth/register", authController.RegisterUser)

	server.HandleFunc("/api/auth/login", authController.LoginUser)
}
