# deployments/ - Deployment Configurations

## Purpose
Contains containerization and orchestration configurations for different deployment environments.

## Structure
```
deployments/
├── docker/              # Docker configurations
└── kubernetes/          # Kubernetes manifests
```

## Docker
- **Dockerfile**: Multi-stage build for Go application
- **docker-compose.yml**: Development environment setup
- **docker-compose.prod.yml**: Production environment setup

## Kubernetes
- **namespace.yaml**: Kubernetes namespace
- **deployment.yaml**: Application deployment
- **service.yaml**: Service definitions
- **configmap.yaml**: Configuration management

## Usage

### Docker
```bash
# Development
docker-compose up -d

# Production
docker-compose -f deployments/docker/docker-compose.prod.yml up -d
```

### Kubernetes
```bash
kubectl apply -f deployments/kubernetes/
```
