package tests

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	resp, err := http.Get("http://localhost:80/health")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestUploadFile(t *testing.T) {
	// Create a test file with unique name
	testContent := "This is a test file for API testing"
	testFile := "test-file-" + t.Name() + ".txt"
	
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFile)

	// Prepare multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	
	fileWriter, err := writer.CreateFormFile("file", testFile)
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	
	_, err = fileWriter.Write([]byte(testContent))
	if err != nil {
		t.Fatalf("Failed to write file content: %v", err)
	}
	
	writer.Close()

	// Make request
	req, err := http.NewRequest("POST", "http://localhost:80/api/v1/files/upload", &buf)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusConflict {
		t.Errorf("Expected status 201 or 409, got %d", resp.StatusCode)
	}
}

func TestListFiles(t *testing.T) {
	resp, err := http.Get("http://localhost:80/api/v1/files")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestStats(t *testing.T) {
	resp, err := http.Get("http://localhost:80/api/v1/stats")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
