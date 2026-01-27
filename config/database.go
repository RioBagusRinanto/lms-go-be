package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database holds the database connection
var DB *gorm.DB

// InitDatabase initializes the database connection
// It establishes a PostgreSQL connection using environment variables
func InitDatabase() *gorm.DB {
	// Get database credentials from environment variables
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "password"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "lms_db"
	}

	// Construct DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = db
	return db
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
