.PHONY: build run test test-coverage docker-build docker-up docker-down migrate-up migrate-down clean

# Build the application
build:
	go build -o bin/labour-thekedar ./cmd/server

# Run the application locally
run:
	go run ./cmd/server

# Run all tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Build docker image
docker-build:
	docker-compose build

# Start docker containers
docker-up:
	docker-compose up -d

# Stop docker containers
docker-down:
	docker-compose down

# Run database migrations up
migrate-up:
	go run ./cmd/server migrate up

# Run database migrations down
migrate-down:
	go run ./cmd/server migrate down

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Tidy go modules
tidy:
	go mod tidy

# Run linter
lint:
	golangci-lint run

# Generate mocks (if using mockery)
mocks:
	mockery --all --dir=internal --output=internal/mocks
