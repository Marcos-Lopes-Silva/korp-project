package config

import (
	"billing/internal/models"
	"os"

	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf("host=%s user=postgres password=postgres dbname=%s port=5432 sslmode=disable", os.Getenv("DB_HOST"), os.Getenv(("DB_NAME")))

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	if err := database.AutoMigrate(&models.Invoice{}); err != nil {
		panic(fmt.Sprintf("Failed to migrate database: %v", err))
	}

	DB = database
}
