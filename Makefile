.PHONY: build run test clean deps

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application  
run:
	go run cmd/server/main.go

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Format code
fmt:
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Generate swagger docs (requires swag)
swagger:
	swag init -g cmd/server/main.go -o ./api/swagger

# Build and run with docker-compose
docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down

# Database operations
migrate:
	./scripts/migrate.sh

seed:
	./scripts/seed.sh
