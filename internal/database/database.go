package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func InitDB() {
	// Load environment variables
	err := godotenv.Load("./../../.env")
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Read configuration
	dbname := os.Getenv("DB_DATABASE")
	password := os.Getenv("DB_PASSWORD")
	username := os.Getenv("DB_USERNAME")
	port := os.Getenv("DB_PORT")
	host := os.Getenv("DB_HOST")

	// Initialize the database connection
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname)
	db, err = sql.Open("mysql", connString)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Verify connection
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established successfully")
}

// GetDB returns the initialized database instance
func GetDB() *sql.DB {
	if db == nil {
		log.Fatal("Database connection is not initialized. Call InitDB first.")
	}
	return db
}
