package main

import (
	"api-file-upload-go/internal/config"
	"api-file-upload-go/internal/database"
	"api-file-upload-go/internal/handlers"
	"api-file-upload-go/internal/logger"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger := logger.New(cfg.LogLevel)

	// Initialize database
	db, err := database.Init(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to initialize database:", err)
	}

	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	r := gin.New()

	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Create file handler
	fileHandler := handlers.NewFileHandler(cfg, db, logger)

	// Setup routes
	handlers.SetupRoutes(r, fileHandler)

	// Get port from config
	port := cfg.Port
	if port == "" {
		logger.Fatal("PORT environment variable is required")
	}

	// Validate port
	if _, err := strconv.Atoi(port); err != nil {
		logger.Fatal("Invalid PORT value:", err)
	}

	// Start server
	logger.Infof("Starting File Upload API on port %s", port)
	if err := r.Run(":" + port); err != nil {
		logger.Fatal("Failed to start server:", err)
	}
}
