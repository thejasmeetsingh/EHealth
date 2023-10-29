package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/thejasmeetsingh/EHealth/internal/database"
)

// Create and return a new DB connection
func getDBConn(isTest bool) *database.Queries {
	dbURL := os.Getenv("DB_URL")

	if isTest {
		dbURL = "postgres://test_db_user:1234@localhost:5432/ehealth_test_db?sslmode=disable"
	}

	if dbURL == "" {
		log.Fatal("DB credentials is not configured in the enviorment")
	}

	conn, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal("Cannot connect to the database: ", err)
	}

	return database.New(conn)
}
