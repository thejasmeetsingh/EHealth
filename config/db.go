package config

import (
	"database/sql"
	"log"
	"os"

	"github.com/thejasmeetsingh/EHealth/internal/database"
)

func GetDBConn() *database.Queries {
	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		log.Fatal("DB credentials is not configured in the enviorment")
	}

	conn, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal("Cannot connect to the database: ", err)
	}

	return database.New(conn)
}
