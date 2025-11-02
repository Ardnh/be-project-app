# ===========================
# Build stage
# ===========================
FROM golang:1.24-alpine AS builder

# Install dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binaries
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/migrate ./cmd/migrate/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/server/main.go


# ===========================
# Final stage
# ===========================
FROM alpine:3.22.2

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binaries from builder
COPY --from=builder /app/migrate /usr/local/bin/migrate
COPY --from=builder /app/main /usr/local/bin/main

# Expose port
EXPOSE 8080

# Run migrations then start server
CMD ["sh", "-c", "migrate && main"]
