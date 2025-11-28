#!/bin/bash
# /opt/deploy/health-check.sh

check_health() {
    local service=$1
    local url=$2

    if curl -f -s "$url" > /dev/null 2>&1; then
        echo "✅ $service is healthy"
        return 0
    else
        echo "❌ $service is unhealthy"
        return 1
    fi
}

echo "🏥 Running health checks..."

check_health "Main App" "http://localhost:8080/health"
check_health "Database" "http://localhost:8080/api/v1/db-health"

# Check container status
echo ""
echo "📦 Container Status:"
docker compose -f /opt/myapp/docker-compose.yml ps

# Check disk space
echo ""
echo "💾 Disk Usage:"
df -h /

# Check memory
echo ""
echo "🧠 Memory Usage:"
free -h
