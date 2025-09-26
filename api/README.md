# api/ - API Documentation

## Purpose
Contains API documentation, specifications, and client tools for the project.

## Structure
```
api/
├── swagger/              # OpenAPI/Swagger documentation
└── postman/             # Postman collections
```

## Contents

### swagger/
- OpenAPI 3.0 specifications
- Auto-generated documentation from code annotations
- Interactive API documentation
- Schema definitions

### postman/
- Postman collection files
- Environment configurations
- Pre-request scripts and tests
- API testing scenarios

## Usage

### Generate Swagger Docs
```bash
make swagger
```

### View Documentation
```bash
# Serve swagger UI locally
swagger-ui-serve api/swagger/swagger.yaml
```

### Import Postman Collection
1. Open Postman
2. Import `api/postman/project_management.postman_collection.json`
3. Set up environment variables
4. Start testing endpoints
