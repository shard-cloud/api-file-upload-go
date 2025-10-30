package database

import (
	"api-file-upload-go/internal/models"
	"fmt"
	"strings"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// convertDatabaseURL converts postgres:// to postgresql:// format and fixes SSL parameters
func convertDatabaseURL(url string) string {
	if strings.HasPrefix(url, "postgres://") {
		url = strings.Replace(url, "postgres://", "postgresql://", 1)
	}
	
	// Replace ssl=true with sslmode=require for PostgreSQL compatibility
	if strings.Contains(url, "ssl=true") {
		url = strings.Replace(url, "ssl=true", "sslmode=require", 1)
	}
	
	return url
}

func Init(databaseURL string) (*gorm.DB, error) {
	// Convert URL format if needed
	convertedURL := convertDatabaseURL(databaseURL)
	
	db, err := gorm.Open(postgres.Open(convertedURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error), // Only log errors
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate
	if err := db.AutoMigrate(&models.File{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}
