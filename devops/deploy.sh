#!/bin/bash
# /opt/deploy/deploy.sh

set -euo pipefail

# ============================================
# Configuration
# ============================================
APP_DIR="/opt/myapp"
BACKUP_DIR="/opt/myapp/backups"
LOG_FILE="/opt/logs/deploy.log"
ROLLBACK_LIMIT=5

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# ============================================
# Logging Functions
# ============================================
log() {
    echo -e "${GREEN}[$(date '+%Y-%m-%d %H:%M:%S')]${NC} $1" | tee -a "$LOG_FILE"
}

error() {
    echo -e "${RED}[$(date '+%Y-%m-%d %H:%M:%S')] ERROR:${NC} $1" | tee -a "$LOG_FILE"
}

warning() {
    echo -e "${YELLOW}[$(date '+%Y-%m-%d %H:%M:%S')] WARNING:${NC} $1" | tee -a "$LOG_FILE"
}

info() {
    echo -e "${BLUE}[$(date '+%Y-%m-%d %H:%M:%S')] INFO:${NC} $1" | tee -a "$LOG_FILE"
}

# ============================================
# Notification Functions
# ============================================
send_notification() {
    local status=$1
    local message=$2

    # Slack notification
    if [ -n "${SLACK_WEBHOOK_URL:-}" ]; then
        curl -X POST "$SLACK_WEBHOOK_URL" \
            -H "Content-Type: application/json" \
            -d "{
                \"text\": \"${status}: ${message}\",
                \"username\": \"Deployment Bot\",
                \"icon_emoji\": \":rocket:\"
            }" 2>/dev/null || true
    fi

    # Email notification (optional)
    # echo "$message" | mail -s "$status" admin@example.com
}

# ============================================
# Backup Functions
# ============================================
backup_database() {
    log "💾 Creating database backup..."

    local backup_name="backup-$(date +%Y%m%d-%H%M%S)"
    mkdir -p "$BACKUP_DIR"

    cd "$APP_DIR"

    # Backup PostgreSQL
    docker compose exec -T postgres pg_dump -U "$DB_USER" "$DB_NAME" | \
        gzip > "$BACKUP_DIR/${backup_name}.sql.gz"

    # Keep only last N backups
    ls -t "$BACKUP_DIR"/*.sql.gz | tail -n +$((ROLLBACK_LIMIT + 1)) | xargs -r rm

    log "✅ Backup created: ${backup_name}.sql.gz"
    echo "$backup_name"
}

backup_containers() {
    log "📦 Backing up container state..."

    local backup_name="containers-$(date +%Y%m%d-%H%M%S)"

    cd "$APP_DIR"

    # Save current image info
    docker compose images --format json > "$BACKUP_DIR/${backup_name}.json"

    log "✅ Container state saved"
}

# ============================================
# Deployment Functions
# ============================================
deploy() {
    log "🚀 Starting deployment..."
    log "Image: ${DEPLOY_IMAGE}"
    log "Branch: ${DEPLOY_BRANCH}"
    log "Commit: ${DEPLOY_COMMIT}"
    log "Deployer: ${DEPLOYER}"

    cd "$APP_DIR" || exit 1

    # Step 1: Create backups
    DB_BACKUP=$(backup_database)
    backup_containers

    # Step 2: Pull latest image
    log "📥 Pulling latest image..."
    if ! docker compose pull; then
        error "Failed to pull image"
        send_notification "❌ FAILED" "Failed to pull image: ${DEPLOY_IMAGE}"
        exit 1
    fi

    # Step 3: Stop current containers
    log "⏹️  Stopping current containers..."
    docker compose down --remove-orphans

    # Step 4: Run database migration
    log "🔄 Running database migration..."
    if ! timeout 300 docker compose run --rm db-migrate; then
        error "Migration failed or timed out!"
        log "🔙 Rolling back..."

        # Restore database
        if [ -f "$BACKUP_DIR/${DB_BACKUP}.sql.gz" ]; then
            warning "Restoring database from backup..."
            gunzip -c "$BACKUP_DIR/${DB_BACKUP}.sql.gz" | \
                docker compose exec -T postgres psql -U "$DB_USER" "$DB_NAME"
        fi

        send_notification "❌ FAILED" "Migration failed. Database restored."
        exit 1
    fi

    log "✅ Migration completed successfully"

    # Step 5: Start new containers
    log "▶️  Starting new containers..."
    if ! docker compose up -d; then
        error "Failed to start containers"
        send_notification "❌ FAILED" "Failed to start containers"
        exit 1
    fi

    # Step 6: Wait for containers to be ready
    log "⏳ Waiting for containers to be ready..."
    sleep 10

    # Step 7: Health checks
    log "🏥 Running health checks..."

    MAX_RETRIES=30
    RETRY_INTERVAL=10

    for i in $(seq 1 $MAX_RETRIES); do
        info "Health check attempt $i/$MAX_RETRIES..."

        # Check if container is running
        if ! docker compose ps main-app | grep -q "Up"; then
            warning "Container is not running"

            if [ $i -eq $MAX_RETRIES ]; then
                error "Container failed to start"
                log "📋 Container logs:"
                docker compose logs --tail=100 main-app

                send_notification "❌ FAILED" "Container failed to start"
                exit 1
            fi

            sleep $RETRY_INTERVAL
            continue
        fi

        # Check health endpoint
        if curl -f -s http://localhost:8080/health > /dev/null 2>&1; then
            log "✅ Application is healthy!"
            break
        fi

        if [ $i -eq $MAX_RETRIES ]; then
            error "Health check failed after $MAX_RETRIES attempts"
            log "📋 Application logs:"
            docker compose logs --tail=100 main-app

            send_notification "❌ FAILED" "Health check failed"
            exit 1
        fi

        sleep $RETRY_INTERVAL
    done

    # Step 8: Smoke tests
    log "🧪 Running smoke tests..."

    # Test critical endpoints
    if ! curl -f -s http://localhost:8080/api/v1/status > /dev/null; then
        error "Status endpoint failed"
        send_notification "❌ FAILED" "Smoke tests failed"
        exit 1
    fi

    # Check version
    VERSION=$(curl -s http://localhost:8080/api/v1/version | jq -r '.version' || echo "unknown")
    log "📊 Deployed version: $VERSION"

    # Step 9: Cleanup
    log "🧹 Cleaning up..."

    # Remove old images (keep last 3)
    docker image prune -af --filter "until=72h" || true

    # Remove old logs
    find /opt/logs -name "*.log" -mtime +30 -delete || true

    # Step 10: Final status
    log "✅ Deployment completed successfully!"
    log "📊 Current status:"
    docker compose ps

    # Send success notification
    send_notification "✅ SUCCESS" "Deployment completed successfully
Environment: Staging
Branch: ${DEPLOY_BRANCH}
Commit: ${DEPLOY_COMMIT}
Version: ${VERSION}
Deployer: ${DEPLOYER}"

    log "🎉 All done!"
}

# ============================================
# Main Execution
# ============================================
main() {
    # Validate required environment variables
    if [ -z "${DEPLOY_IMAGE:-}" ]; then
        error "DEPLOY_IMAGE not set"
        exit 1
    fi

    # Load environment
    if [ -f "$APP_DIR/.env" ]; then
        set -a
        source "$APP_DIR/.env"
        set +a
    fi

    # Run deployment
    if deploy; then
        exit 0
    else
        error "Deployment failed"
        exit 1
    fi
}

# Trap errors
trap 'error "Deployment failed with error on line $LINENO"' ERR

main "$@"
