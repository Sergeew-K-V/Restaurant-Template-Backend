package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func NewDBConfig() *DBConfig {
	return &DBConfig{
		Host:     "localhost",
		Port:     8888,
		User:     "admin",
		Password: "admin",
		DBName:   "go-postgres",
	}
}

func GetDBConnection(c *DBConfig) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DBName)
	fmt.Println("POSTGRESSQL INFO:", psqlInfo)

	db, err := sql.Open("postgres", psqlInfo)
	fmt.Println("GetDBConnection db:", db)

	if err != nil {
		return nil, fmt.Errorf("error opening connection to DB: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to DB: %v", err)
	}

	log.Println("Successfully connected to the PostgreSQL database")
	return db, nil
}

func CloseDB(db *sql.DB) {
	if db != nil {
		db.Close()
		log.Println("The database connection was closed.")
	}
}
