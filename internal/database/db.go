package database

import (
	"log"

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

	log.Println("Database connected and migrated successfully")
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
