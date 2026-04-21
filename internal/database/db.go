package database

import (
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"myplants-server/internal/models"
)

var DB *gorm.DB

// InitDB initializes the database connection and performs auto-migration
func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("myplants.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the schema
	err = DB.AutoMigrate(&models.Content{}, &models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Optimize connection pool for 1-core / 1GB VPS
	if sqlDB, err := DB.DB(); err == nil {
		sqlDB.SetMaxOpenConns(1)
		sqlDB.SetMaxIdleConns(1)
		sqlDB.SetConnMaxLifetime(5 * time.Minute)
	}

	log.Println("Database connected and migrated successfully")
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
