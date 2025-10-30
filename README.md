# File Upload API Go

REST API for uploading files to a storage system and tracking metadata in a database. Built with Go, Gin framework, GORM, and PostgreSQL.

## 🎯 Features

- ✅ Upload files with validation via REST API
- ✅ List uploaded files with pagination
- ✅ Download files via API
- ✅ Delete files via API
- ✅ Show upload statistics
- ✅ File deduplication by hash
- ✅ PostgreSQL with GORM ORM
- ✅ Docker support
- ✅ Structured logging
- ✅ Environment configuration
- ✅ CORS support
- ✅ Health check endpoint

## 📋 Requirements

- Go 1.21+
- PostgreSQL 14+ (or use Docker Compose)
- Docker (optional)

## 🚀 How to run

### Linux/macOS
```bash
# Copy .env and configure DATABASE
cp env.example .env

# Install dependencies
go mod tidy

# Build binary
go build -o api-file-upload-go ./cmd/main.go

# Run API
./api-file-upload-go
```

### Windows
```bash
# Copy .env and configure DATABASE
copy env.example .env

# Install dependencies
go mod tidy

# Build binary
go build -o api-file-upload-go.exe ./cmd/main.go

# Run API
.\api-file-upload-go.exe
```

## 📦 API Endpoints

### File Operations
- `POST /api/v1/files/upload` – Upload a file
- `GET /api/v1/files` – List uploaded files (with pagination)
- `GET /api/v1/files/:id` – Get file details by ID
- `GET /api/v1/files/:id/download` – Download file by ID
- `DELETE /api/v1/files/:id` – Delete file by ID

### Statistics
- `GET /api/v1/stats` – Show upload statistics

### System
- `GET /health` – Health check endpoint
- `GET /` – API information

## 🗄️ Database

Connection via `DATABASE` environment variable:

```env
DATABASE=postgres://user:pass@localhost:5432/uploader_db
```

## 📂 Structure

```
cli-uploader-go/
├── cmd/                    # Main application entry point
├── internal/
│   ├── handlers/           # HTTP handlers (upload, list, download, delete, stats)
│   ├── config/            # Configuration management
│   ├── database/          # Database connection and models
│   ├── logger/            # Structured logging
│   ├── models/            # Data models
│   └── utils/             # Utility functions
├── docs/                  # Complete documentation
└── tests/                 # Test files
```

## 🔧 Configuration

### Environment Variables

```env
# Database (required)
DATABASE=postgres://user:pass@localhost:5432/uploader_db

# Upload configuration
UPLOAD_DIR=./uploads
MAX_FILE_SIZE=10485760
ALLOWED_EXTENSIONS=.jpg,.jpeg,.png,.gif,.pdf,.txt,.doc,.docx

# Server configuration
PORT=80
ENVIRONMENT=development

# Logging
LOG_LEVEL=info
```

## 🐳 Docker

### Build and run
```bash
# Build image
docker build -t api-file-upload-go .

# Run with environment variables
docker run --rm -p 80:80 \
  -e DATABASE=postgres://user:pass@host:5432/db \
  -e PORT=80 \
  -e UPLOAD_DIR=/app/uploads \
  -v $(pwd)/uploads:/app/uploads \
  api-file-upload-go
```

### Docker Compose
```bash
# Start services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## 📊 Usage Examples

### Upload a file
```bash
curl -X POST -F "file=@document.pdf" http://localhost:80/api/v1/files/upload
```

### List files
```bash
curl http://localhost:80/api/v1/files?limit=10&offset=0
```

### Download a file
```bash
curl http://localhost:80/api/v1/files/1/download -o downloaded_file.pdf
```

### Get file details
```bash
curl http://localhost:80/api/v1/files/1
```

### Show statistics
```bash
curl http://localhost:80/api/v1/stats
```

### Delete a file
```bash
curl -X DELETE http://localhost:80/api/v1/files/1
```

### Health check
```bash
curl http://localhost:80/health
```

## 🧪 Testing

### Unit tests
```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

### API tests
```bash
# Start the API first
./api-file-upload-go  # Linux/macOS
# or
.\api-file-upload-go.exe  # Windows

# Run API tests (in another terminal)
go test ./tests/
```

## 📄 License

MIT