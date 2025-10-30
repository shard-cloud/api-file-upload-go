package utils

import (
	"mime"
	"path/filepath"
	"strings"
)

// Contains checks if a slice contains a string
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// GetMimeType returns the MIME type of a file
func GetMimeType(filePath string) string {
	ext := filepath.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		return "application/octet-stream"
	}
	return strings.Split(mimeType, ";")[0]
}
