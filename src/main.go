package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"restaurant-backend/src/config"
	"restaurant-backend/src/database"
	"restaurant-backend/src/services/auth"
	"restaurant-backend/src/services/router"
	"restaurant-backend/src/types"
	"strconv"

	"github.com/rs/cors"
)

func main() {
	envConfig := config.LoadGlobalConfig()

	db, err := database.GetDBConnection(envConfig.DB)

	if err != nil {
		log.Fatal("Error connecting to DB:", err)
	}
	defer database.CloseDB(db)

	mux := http.NewServeMux()
	fmt.Printf("Server started on %d port \n", envConfig.App.Port)

	// CORS configuration using github.com/rs/cors
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	server := c.Handler(mux)

	ctx := context.Background()

	ctx = context.WithValue(ctx, types.CookieSecretKey, envConfig.App.CookieSecretKey)

	router, err := router.InitRouter(ctx, mux, "/api", "")
	if err != nil {
		log.Fatal("Error initializing router:", err)
	}

	authService, err := auth.InitService(ctx, db)
	if err != nil {
		log.Fatal("Error initializing auth service:", err)
	}

	router.RegistryPublicRoute(ctx, "/auth/register", func(w http.ResponseWriter, r *http.Request) {
		authService.RegisterUser(ctx, w, r)
	})
	router.RegistryPublicRoute(ctx, "/auth/login", func(w http.ResponseWriter, r *http.Request) {
		authService.LoginUser(ctx, w, r)
	})
	router.RegistryPrivateRoute(ctx, "/me", func(w http.ResponseWriter, r *http.Request) {
		authService.Me(ctx, w, r)
	})

	portStr := ":" + strconv.Itoa(envConfig.App.Port)

	log.Fatal(http.ListenAndServe(portStr, server))
}
