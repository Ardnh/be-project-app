# scripts/ - Build and Deployment Scripts

## Purpose
Contains automation scripts for building, testing, and deploying the application.

## Scripts

### build.sh
- Build the Go application
- Set build information (version, commit, build time)
- Create binaries for different platforms

### migrate.sh  
- Run database migrations
- Wait for database connectivity
- Handle migration rollbacks

### deploy.sh
- Deploy to different environments
- Build and push Docker images
- Update Kubernetes deployments

### seed.sh
- Seed database with initial data
- Create default users and categories
- Set up development data

## Usage
```bash
# Make scripts executable
chmod +x scripts/*.sh

# Run specific script
./scripts/build.sh
./scripts/migrate.sh
```
