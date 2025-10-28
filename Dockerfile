# Build stage
FROM golang:1.24-alpine AS builder

# Install dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /bin/migrate ./cmd/migrate/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /bin/main ./cmd/server/main.go

# Final stage
FROM alpine:3.22.2

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Run the application
CMD ["sh", "-c", "./bin/migrate/migrate && ./bin/main/main"]
