package config

import (
	"billing/internal/models"

	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=localhost user=postgres password=postgres dbname=korp_invoice port=5432 sslmode=disable"

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	if err := database.AutoMigrate(&models.Invoice{}); err != nil {
		panic(fmt.Sprintf("Failed to migrate database: %v", err))
	}

	DB = database
}
