package app

import (
	"database/sql"
	"fmt"
	"time"

	// Import mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// DB is the global database connection
var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
	// Default database configuration
	host := "localhost"
	port := 3306
	username := "root"
	password := ""
	dbname := "goapp"

	// Build connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username,
		password,
		host,
		port,
		dbname,
	)

	// Try to connect to the database
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		panic(err)
	}

	// Set connection pool settings
	DB.SetMaxIdleConns(10)
	DB.SetMaxOpenConns(100)
	DB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := DB.Ping(); err != nil {
		fmt.Printf("Failed to ping database: %v\n", err)
		panic(err)
	}

	fmt.Println("Database connected successfully")
	Info("Database connected successfully")
}

// Query is a helper function to execute a query and scan the results
func Query(query string, args ...interface{}) (*sql.Rows, error) {
	return DB.Query(query, args...)
}

// Exec is a helper function to execute a statement
func Exec(query string, args ...interface{}) (sql.Result, error) {
	return DB.Exec(query, args...)
}

// QueryRow is a helper function to execute a query that returns a single row
func QueryRow(query string, args ...interface{}) *sql.Row {
	return DB.QueryRow(query, args...)
}

// Begin starts a transaction
func Begin() (*sql.Tx, error) {
	return DB.Begin()
}
