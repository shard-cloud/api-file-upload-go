package handlers

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(r *gin.Engine, fileHandler *FileHandler) {
	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// File routes
		files := v1.Group("/files")
		{
			files.POST("/upload", fileHandler.UploadFile)
			files.GET("", fileHandler.ListFiles)
			files.GET("/:id", fileHandler.GetFile)
			files.GET("/:id/download", fileHandler.DownloadFile)
			files.DELETE("/:id", fileHandler.DeleteFile)
		}

		// Stats route
		v1.GET("/stats", fileHandler.GetStats)
	}

	// Health check route
	r.GET("/health", fileHandler.HealthCheck)

	// Root route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "File Upload API",
			"version": "1.0.0",
			"docs":    "/api/v1",
			"health":  "/health",
		})
	})
}
