FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/migrate ./cmd/migrate/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/server/main.go

FROM alpine:3.22.2
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/migrate /usr/local/bin/migrate
COPY --from=builder /app/main /usr/local/bin/main
EXPOSE 8080
CMD ["sh", "-c", "migrate && main"]
