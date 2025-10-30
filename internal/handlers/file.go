package handlers

import (
	"api-file-upload-go/internal/config"
	"api-file-upload-go/internal/models"
	"api-file-upload-go/internal/utils"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type FileHandler struct {
	config *config.Config
	db     *gorm.DB
	logger *logrus.Logger
}

func NewFileHandler(cfg *config.Config, db *gorm.DB, logger *logrus.Logger) *FileHandler {
	return &FileHandler{
		config: cfg,
		db:     db,
		logger: logger,
	}
}

// UploadFile handles file upload
func (h *FileHandler) UploadFile(c *gin.Context) {
	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "No file uploaded",
		})
		return
	}

	// Check file size (only if MaxFileSize is defined)
	if h.config.MaxFileSize > 0 && file.Size > h.config.MaxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": fmt.Sprintf("File size exceeds maximum allowed size: %d bytes", h.config.MaxFileSize),
		})
		return
	}

	// Get file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))

	// Check file extension (only if AllowedExtensions is defined)
	if len(h.config.AllowedExtensions) > 0 {
		if !utils.Contains(h.config.AllowedExtensions, ext) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": fmt.Sprintf("File extension not allowed: %s", ext),
			})
			return
		}
	}

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(h.config.UploadDir, 0755); err != nil {
		h.logger.Error("Failed to create upload directory:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to create upload directory",
		})
		return
	}

	// Generate unique filename
	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	destPath := filepath.Join(h.config.UploadDir, fileName)

	// Save uploaded file
	if err := c.SaveUploadedFile(file, destPath); err != nil {
		h.logger.Error("Failed to save uploaded file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to save uploaded file",
		})
		return
	}

	// Calculate file hash
	hash, err := h.calculateFileHash(destPath)
	if err != nil {
		// Clean up uploaded file if hash calculation fails
		os.Remove(destPath)
		h.logger.Error("Failed to calculate file hash:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to calculate file hash",
		})
		return
	}

	// Check if file already exists
	var existingFile models.File
	if err := h.db.Where("hash = ?", hash).First(&existingFile).Error; err == nil {
		// Clean up uploaded file if duplicate
		os.Remove(destPath)
		c.JSON(http.StatusConflict, gin.H{
			"error":   true,
			"message": "File already exists",
			"file_id": existingFile.ID,
		})
		return
	}

	// Get MIME type
	mimeType := utils.GetMimeType(file.Filename)

	// Create file record
	fileRecord := models.File{
		Name:         fileName,
		OriginalName: file.Filename,
		Path:         destPath,
		Size:         file.Size,
		MimeType:     mimeType,
		Extension:    ext,
		Hash:         hash,
	}

	// Save to database
	if err := h.db.Create(&fileRecord).Error; err != nil {
		// Clean up uploaded file if database save fails
		os.Remove(destPath)
		h.logger.Error("Failed to save file metadata:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to save file metadata",
		})
		return
	}

	h.logger.Infof("File uploaded successfully: %s (ID: %d)", fileRecord.OriginalName, fileRecord.ID)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "File uploaded successfully",
		"file": gin.H{
			"id":           fileRecord.ID,
			"name":         fileRecord.OriginalName,
			"size":         fileRecord.Size,
			"mime_type":    fileRecord.MimeType,
			"extension":    fileRecord.Extension,
			"hash":         fileRecord.Hash,
			"uploaded_at":  fileRecord.UploadedAt,
		},
	})
}

// ListFiles handles file listing with pagination
func (h *FileHandler) ListFiles(c *gin.Context) {
	// Parse pagination parameters
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	var files []models.File

	// Get total count
	var total int64
	if err := h.db.Model(&models.File{}).Count(&total).Error; err != nil {
		h.logger.Error("Failed to count files:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to count files",
		})
		return
	}

	// Get files with pagination
	if err := h.db.Limit(limit).Offset(offset).Order("uploaded_at DESC").Find(&files).Error; err != nil {
		h.logger.Error("Failed to list files:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to list files",
		})
		return
	}

	// Format response
	var fileList []gin.H
	for _, file := range files {
		fileList = append(fileList, gin.H{
			"id":           file.ID,
			"name":         file.OriginalName,
			"size":         file.Size,
			"mime_type":    file.MimeType,
			"extension":    file.Extension,
			"hash":         file.Hash,
			"uploaded_at":  file.UploadedAt,
			"download_url": fmt.Sprintf("/api/v1/files/%d/download", file.ID),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"files":  fileList,
			"total":  total,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// GetFile handles getting file details by ID
func (h *FileHandler) GetFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid file ID",
		})
		return
	}

	var file models.File
	if err := h.db.First(&file, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   true,
				"message": "File not found",
			})
			return
		}
		h.logger.Error("Failed to get file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to get file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id":           file.ID,
			"name":         file.OriginalName,
			"size":         file.Size,
			"mime_type":    file.MimeType,
			"extension":    file.Extension,
			"hash":         file.Hash,
			"uploaded_at":  file.UploadedAt,
			"download_url": fmt.Sprintf("/api/v1/files/%d/download", file.ID),
		},
	})
}

// DownloadFile handles file download
func (h *FileHandler) DownloadFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid file ID",
		})
		return
	}

	var file models.File
	if err := h.db.First(&file, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   true,
				"message": "File not found",
			})
			return
		}
		h.logger.Error("Failed to get file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to get file",
		})
		return
	}

	// Check if file exists on disk
	if _, err := os.Stat(file.Path); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "File not found on disk",
		})
		return
	}

	// Set headers for file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.OriginalName))
	c.Header("Content-Type", file.MimeType)

	// Serve file
	c.File(file.Path)
}

// DeleteFile handles file deletion
func (h *FileHandler) DeleteFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid file ID",
		})
		return
	}

	var file models.File
	if err := h.db.First(&file, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   true,
				"message": "File not found",
			})
			return
		}
		h.logger.Error("Failed to get file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to get file",
		})
		return
	}

	// Delete file from disk
	if err := os.Remove(file.Path); err != nil {
		h.logger.Warn("Failed to delete file from disk:", err)
	}

	// Delete from database (soft delete)
	if err := h.db.Delete(&file).Error; err != nil {
		h.logger.Error("Failed to delete file from database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to delete file from database",
		})
		return
	}

	h.logger.Infof("File deleted successfully: %s (ID: %d)", file.OriginalName, file.ID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "File deleted successfully",
	})
}

// GetStats handles upload statistics
func (h *FileHandler) GetStats(c *gin.Context) {
	// Get total count
	var totalFiles int64
	if err := h.db.Model(&models.File{}).Count(&totalFiles).Error; err != nil {
		h.logger.Error("Failed to count files:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to count files",
		})
		return
	}

	// Get total size
	var totalSize int64
	if err := h.db.Model(&models.File{}).Select("COALESCE(SUM(size), 0)").Scan(&totalSize).Error; err != nil {
		h.logger.Error("Failed to calculate total size:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to calculate total size",
		})
		return
	}

	// Get files by extension
	var extensionStats []struct {
		Extension string `json:"extension"`
		Count     int64  `json:"count"`
		Size      int64  `json:"size"`
	}
	if err := h.db.Model(&models.File{}).
		Select("extension, COUNT(*) as count, SUM(size) as size").
		Group("extension").
		Order("count DESC").
		Scan(&extensionStats).Error; err != nil {
		h.logger.Error("Failed to get extension stats:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to get extension stats",
		})
		return
	}

	// Get recent uploads (last 24 hours)
	var recentUploads int64
	yesterday := time.Now().Add(-24 * time.Hour)
	if err := h.db.Model(&models.File{}).
		Where("uploaded_at > ?", yesterday).
		Count(&recentUploads).Error; err != nil {
		h.logger.Error("Failed to count recent uploads:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to count recent uploads",
		})
		return
	}

	// Get largest file
	var largestFile models.File
	if err := h.db.Order("size DESC").First(&largestFile).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			h.logger.Error("Failed to get largest file:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Failed to get largest file",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"total_files":      totalFiles,
			"total_size":       totalSize,
			"recent_uploads":   recentUploads,
			"largest_file": gin.H{
				"name": largestFile.OriginalName,
				"size": largestFile.Size,
			},
			"extension_stats": extensionStats,
		},
	})
}

// HealthCheck handles health check endpoint
func (h *FileHandler) HealthCheck(c *gin.Context) {
	// Test database connection
	if err := h.db.Raw("SELECT 1").Error; err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": "Database connection failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "API is healthy",
		"version": "1.0.0",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// calculateFileHash calculates MD5 hash of a file
func (h *FileHandler) calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
