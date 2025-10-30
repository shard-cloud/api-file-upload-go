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
		Database:          os.Getenv("DATABASE"),
		Port:              os.Getenv("PORT"),
		UploadDir:         os.Getenv("UPLOAD_DIR"),
		MaxFileSize:       maxFileSize,
		AllowedExtensions: allowedExtensions,
		LogLevel:          os.Getenv("LOG_LEVEL"),
		Environment:       os.Getenv("ENVIRONMENT"),
	}
}
