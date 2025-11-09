package app

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// DB is the global database connection
var DB *gorm.DB

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

	// Configure GORM
	config := &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	}

	// Try to connect to the database
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		panic(err)
	}

	// Set connection pool settings
	sqlDB, err := DB.DB()
	if err != nil {
		fmt.Printf("Failed to get underlying *sql.DB: %v\n", err)
		panic(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Println("Database connected successfully")
	Info("Database connected successfully")
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// WithTx executes function within transaction
func WithTx(fn func(tx *gorm.DB) error) error {
	return DB.Transaction(fn)
}
