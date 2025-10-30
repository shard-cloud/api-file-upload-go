FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o api-file-upload-go cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates curl
WORKDIR /root/

COPY --from=builder /app/api-file-upload-go .

# Create uploads directory
RUN mkdir -p /app/uploads

# Expose port 80
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 \
  CMD curl --fail http://localhost:80/health || exit 1

CMD ["./api-file-upload-go"]
