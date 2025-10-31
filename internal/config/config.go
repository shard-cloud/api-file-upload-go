package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Database          string
	Port              string
	UploadDir         string
	MaxFileSize       int64
	AllowedExtensions []string
	LogLevel          string
	Environment       string
}

func Load() *Config {
	maxFileSize := int64(0)
	if maxSizeStr := os.Getenv("MAX_FILE_SIZE"); maxSizeStr != "" {
		if parsed, err := strconv.ParseInt(maxSizeStr, 10, 64); err == nil {
			maxFileSize = parsed
		}
	}

	allowedExtensions := []string{}
	if extStr := os.Getenv("ALLOWED_EXTENSIONS"); extStr != "" {
		allowedExtensions = strings.Split(extStr, ",")
	}

	return &Config{
		Database:          normalizeDatabaseURL(os.Getenv("DATABASE")),
		Port:              os.Getenv("PORT"),
		UploadDir:         os.Getenv("UPLOAD_DIR"),
		MaxFileSize:       maxFileSize,
		AllowedExtensions: allowedExtensions,
		LogLevel:          os.Getenv("LOG_LEVEL"),
		Environment:       os.Getenv("ENVIRONMENT"),
	}
}

// normalizeDatabaseURL fixes common SSL parameter issues in PostgreSQL connection strings
func normalizeDatabaseURL(url string) string {
	if url == "" {
		return url
	}

	// Replace ssl=true with sslmode=require (PostgreSQL standard)
	url = strings.ReplaceAll(url, "ssl=true", "sslmode=require")
	url = strings.ReplaceAll(url, "ssl=false", "sslmode=disable")
	
	// If no sslmode is specified and it's a remote database, add sslmode=require
	if !strings.Contains(url, "sslmode=") && 
	   !strings.Contains(url, "localhost") && 
	   !strings.Contains(url, "127.0.0.1") {
		if strings.Contains(url, "?") {
			url += "&sslmode=require"
		} else {
			url += "?sslmode=require"
		}
	}

	return url
}
