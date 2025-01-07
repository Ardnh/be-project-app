# Dockerfile

# Stage 1: Build
FROM golang:1.20-alpine AS builder

# Menetapkan direktori kerja di dalam container
WORKDIR /app

# Menyalin go.mod dan go.sum terlebih dahulu untuk caching dependensi
COPY go.mod go.sum ./

# Mengunduh dependensi
RUN go mod download

# Menyalin seluruh kode aplikasi
COPY . .

# Membuat binary statis
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 2: Run
FROM alpine:latest

# Menetapkan direktori kerja di dalam container
WORKDIR /root/

# Menyalin binary dari stage sebelumnya
COPY --from=builder /app/main .

# Menjalankan aplikasi
CMD ["./main"]
