package main

import (
	"fmt"
	"log"
	"net/http"
	"restaurant-backend/src/config"
	"restaurant-backend/src/database"
	"restaurant-backend/src/models"
	"restaurant-backend/src/routes"
	"strconv"

	"github.com/rs/cors"
)

func main() {
	var AppContext models.AppContext
	envConfig := config.LoadGlobalConfig()

	db, err := database.GetDBConnection(envConfig.DB)

	AppContext.Config = envConfig
	AppContext.DB = db
	if err != nil {
		log.Fatal("Error connecting to DB:", err)
	}
	defer database.CloseDB(db)

	mux := http.NewServeMux()
	AppContext.Mux = mux
	fmt.Printf("Server started on %d port \n", envConfig.App.Port)

	routes.AuthRoutes(&AppContext)

	// CORS configuration using github.com/rs/cors
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	server := c.Handler(mux)
	portStr := ":" + strconv.Itoa(envConfig.App.Port)

	log.Fatal(http.ListenAndServe(portStr, server))
}
