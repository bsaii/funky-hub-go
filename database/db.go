package database

import (
	"database/sql"
	"log"
)

var DB *sql.DB

func InitDb(connStr string) {
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	pingErr := DB.Ping()
	if pingErr != nil {
		log.Fatalf("Failed to ping the database %v", pingErr)
	}

	log.Println("Connected to the database...")
}
