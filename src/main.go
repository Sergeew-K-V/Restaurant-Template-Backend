package main

import (
	"fmt"
	"log"
	"net/http"
	"restaurant-backend/src/config"
	"restaurant-backend/src/controllers"
	"restaurant-backend/src/database"
	"strconv"
)

func main() {
	cfg := config.LoadGlobalConfig()

	fmt.Printf("Server started on %d port \n", cfg.App.Port)
	portStr := ":" + strconv.Itoa(cfg.App.Port)

	http.HandleFunc("/", controllers.RegisterUser)

	config := database.NewDBConfig()

	db, err := database.GetDBConnection(config)

	if err != nil {
		log.Fatal("Error connecting to DB:", err)
	}
	defer database.CloseDB(db)

	fmt.Printf("db %v \n", db)

	log.Fatal(http.ListenAndServe(portStr, nil))
}
