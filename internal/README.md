# internal/ - Private Application Code

## Purpose

Contains the private application and library code. Code in this directory cannot be imported by external applications.

## Architecture Layers

### 🏗️ Clean Architecture Implementation

```
┌─────────────────────────────────────────────────────────┐
│                    External Interfaces                  │
│                   (HTTP, gRPC, CLI)                     │
├─────────────────────────────────────────────────────────┤
│                   Application Layer                     │
│              (Use Cases, DTOs, Services)                │
├─────────────────────────────────────────────────────────┤
│                     Domain Layer                        │
│           (Entities, Repositories, Services)            │
├─────────────────────────────────────────────────────────┤
│                 Infrastructure Layer                    │
│            (Database, Cache, External APIs)             │
└─────────────────────────────────────────────────────────┘
```

## Directory Structure

### `/domain` - Business Logic Core

- **Entities**: Core business objects
- **Repositories**: Data access interfaces
- **Services**: Business logic interfaces

### `/application` - Application Services

- **Services**: Use case implementations
- **DTOs**: Data Transfer Objects

### `/infrastructure` - External Dependencies

- **Database**: Database implementations
- **Cache**: Caching implementations
- **External**: Third-party integrations

### `/interfaces` - External Communication

- **HTTP**: REST API handlers
- **gRPC**: gRPC service implementations

### `/utils` - Helper Functions

- Common utilities and helpers
- Cross-cutting concerns

### `/config` - Configuration

- Configuration structures
- Environment setup

## Design Principles

1. **Dependency Inversion**: Depend on interfaces, not implementations
2. **Single Responsibility**: Each package has one reason to change
3. **Open/Closed**: Open for extension, closed for modification
4. **Interface Segregation**: Small, focused interfaces
5. **Dependency Injection**: Wire dependencies at startup

## Key Rules

- Domain layer **never** depends on infrastructure
- Infrastructure **implements** domain interfaces
- Application **orchestrates** domain operations
- Interfaces **adapt** external communications
