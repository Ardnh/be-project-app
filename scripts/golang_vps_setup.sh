#!/bin/bash

# ==========================================
# Setup Deployment Golang Project ke VPS
# Stack: Golang, PostgreSQL, Redis, Docker
# ==========================================

set -e  # Exit on error

# Warna untuk output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Setup Deployment Golang Project${NC}"
echo -e "${GREEN}========================================${NC}"

# ==========================================
# 1. Update sistem
# ==========================================
echo -e "\n${YELLOW}[1/7] Updating system...${NC}"
sudo apt update && sudo apt upgrade -y

# ==========================================
# 2. Install Docker
# ==========================================
echo -e "\n${YELLOW}[2/7] Installing Docker...${NC}"
if ! command -v docker &> /dev/null; then
    sudo apt install -y apt-transport-https ca-certificates curl software-properties-common
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    sudo apt update
    sudo apt install -y docker-ce docker-ce-cli containerd.io
    sudo systemctl start docker
    sudo systemctl enable docker
    sudo usermod -aG docker $USER
    echo -e "${GREEN}Docker installed successfully${NC}"
else
    echo -e "${GREEN}Docker already installed${NC}"
fi

# ==========================================
# 3. Install Docker Compose
# ==========================================
echo -e "\n${YELLOW}[3/7] Installing Docker Compose...${NC}"
if ! command -v docker-compose &> /dev/null; then
    sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
    echo -e "${GREEN}Docker Compose installed successfully${NC}"
else
    echo -e "${GREEN}Docker Compose already installed${NC}"
fi

# ==========================================
# 4. Install Git
# ==========================================
echo -e "\n${YELLOW}[4/7] Installing Git...${NC}"
if ! command -v git &> /dev/null; then
    sudo apt install -y git
    echo -e "${GREEN}Git installed successfully${NC}"
else
    echo -e "${GREEN}Git already installed${NC}"
fi

# ==========================================
# 5. Setup Project Directory
# ==========================================
echo -e "\n${YELLOW}[5/7] Setting up project directory...${NC}"
PROJECT_DIR="/var/www/golang-app"
sudo mkdir -p $PROJECT_DIR
sudo chown -R $USER:$USER $PROJECT_DIR
echo -e "${GREEN}Project directory created at $PROJECT_DIR${NC}"

# ==========================================
# 6. Create Docker Compose file
# ==========================================
echo -e "\n${YELLOW}[6/7] Creating Docker Compose configuration...${NC}"
cat > $PROJECT_DIR/docker-compose.yml << 'EOF'
version: '3.8'

services:
  # PostgreSQL Database
  postgres:
    image: postgres:15-alpine
    container_name: golang_postgres
    restart: always
    environment:
      POSTGRES_DB: ${DB_NAME:-myapp}
      POSTGRES_USER: ${DB_USER:-postgres}
      POSTGRES_PASSWORD: ${DB_PASSWORD:-postgres123}
    ports:
      - "${DB_PORT:-5432}:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - app_network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER:-postgres}"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Redis Cache
  redis:
    image: redis:7-alpine
    container_name: golang_redis
    restart: always
    command: redis-server --requirepass ${REDIS_PASSWORD:-redis123}
    ports:
      - "${REDIS_PORT:-6379}:6379"
    volumes:
      - redis_data:/data
    networks:
      - app_network
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Golang Application
  app:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: golang_app
    restart: always
    ports:
      - "${APP_PORT:-8080}:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=${DB_USER:-postgres}
      - DB_PASSWORD=${DB_PASSWORD:-postgres123}
      - DB_NAME=${DB_NAME:-myapp}
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - REDIS_PASSWORD=${REDIS_PASSWORD:-redis123}
      - APP_ENV=${APP_ENV:-production}
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - app_network

volumes:
  postgres_data:
  redis_data:

networks:
  app_network:
    driver: bridge
EOF

echo -e "${GREEN}Docker Compose file created${NC}"

# ==========================================
# 7. Create Example Dockerfile
# ==========================================
echo -e "\n${YELLOW}[7/7] Creating example Dockerfile...${NC}"
cat > $PROJECT_DIR/Dockerfile << 'EOF'
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
EOF

echo -e "${GREEN}Dockerfile created${NC}"

# ==========================================
# Create .env.example file
# ==========================================
cat > $PROJECT_DIR/.env.example << 'EOF'
# Database Configuration
DB_NAME=myapp
DB_USER=postgres
DB_PASSWORD=postgres123
DB_PORT=5432

# Redis Configuration
REDIS_PASSWORD=redis123
REDIS_PORT=6379

# Application Configuration
APP_PORT=8080
APP_ENV=production
EOF

echo -e "${GREEN}.env.example file created${NC}"

# ==========================================
# Create deployment script
# ==========================================
cat > $PROJECT_DIR/deploy.sh << 'EOF'
#!/bin/bash

set -e

echo "Starting deployment..."

# Pull latest code
git pull origin main

# Stop existing containers
docker-compose down

# Build and start containers
docker-compose up -d --build

# Show logs
docker-compose logs -f
EOF

chmod +x $PROJECT_DIR/deploy.sh
echo -e "${GREEN}Deployment script created${NC}"

# ==========================================
# Setup Firewall (UFW)
# ==========================================
echo -e "\n${YELLOW}Setting up firewall...${NC}"
if command -v ufw &> /dev/null; then
    sudo ufw allow 22/tcp
    sudo ufw allow 80/tcp
    sudo ufw allow 443/tcp
    sudo ufw allow 8080/tcp
    sudo ufw --force enable
    echo -e "${GREEN}Firewall configured${NC}"
fi

# ==========================================
# Final Instructions
# ==========================================
echo -e "\n${GREEN}========================================${NC}"
echo -e "${GREEN}Setup completed successfully!${NC}"
echo -e "${GREEN}========================================${NC}"
echo -e "\n${YELLOW}Next steps:${NC}"
echo -e "1. Copy your Golang project to: ${GREEN}$PROJECT_DIR${NC}"
echo -e "2. Create .env file: ${GREEN}cp $PROJECT_DIR/.env.example $PROJECT_DIR/.env${NC}"
echo -e "3. Edit .env file with your credentials"
echo -e "4. Navigate to project: ${GREEN}cd $PROJECT_DIR${NC}"
echo -e "5. Start services: ${GREEN}docker-compose up -d${NC}"
echo -e "6. View logs: ${GREEN}docker-compose logs -f${NC}"
echo -e "7. Stop services: ${GREEN}docker-compose down${NC}"
echo -e "\n${YELLOW}Useful commands:${NC}"
echo -e "- Check status: ${GREEN}docker-compose ps${NC}"
echo -e "- Restart app: ${GREEN}docker-compose restart app${NC}"
echo -e "- View app logs: ${GREEN}docker-compose logs -f app${NC}"
echo -e "- Access PostgreSQL: ${GREEN}docker exec -it golang_postgres psql -U postgres -d myapp${NC}"
echo -e "- Access Redis CLI: ${GREEN}docker exec -it golang_redis redis-cli -a redis123${NC}"
echo -e "\n${RED}IMPORTANT: Remember to logout and login again for Docker group changes to take effect!${NC}"
