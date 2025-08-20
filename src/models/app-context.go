package models

import (
	"database/sql"
	"net/http"
	"restaurant-backend/src/config"
)

type AppContext struct {
	DB     *sql.DB
	Mux    *http.ServeMux
	Config *config.GlobalConfig
}
