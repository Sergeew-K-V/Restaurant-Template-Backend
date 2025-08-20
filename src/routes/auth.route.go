package routes

import (
	"restaurant-backend/src/controllers"
	"restaurant-backend/src/models"
)

func AuthRoutes(context *models.AppContext) {
	authController := controllers.NewAuthController(context)

	context.Mux.HandleFunc("/api/auth/register", authController.RegisterUser)

	context.Mux.HandleFunc("/api/auth/login", authController.LoginUser)
}
