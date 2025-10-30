.PHONY: dev test lint format clean build build-linux build-windows docker-build docker-run

dev:
	go run cmd/main.go

test:
	go test ./...

lint:
	golangci-lint run

format:
	gofmt -s -w .
	goimports -w .

clean:
	go clean
	rm -f api-file-upload-go api-file-upload-go.exe

build:
	go build -o api-file-upload-go cmd/main.go

build-linux:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api-file-upload-go cmd/main.go

build-windows:
	CGO_ENABLED=0 GOOS=windows go build -a -installsuffix cgo -o api-file-upload-go.exe cmd/main.go

docker-build:
	docker build -t api-file-upload-go .

docker-run:
	docker run -p 80:80 \
		-e DATABASE="$(DATABASE)" \
		-e PORT=80 \
		-e UPLOAD_DIR="/app/uploads" \
		-v $(PWD)/uploads:/app/uploads \
		api-file-upload-go
